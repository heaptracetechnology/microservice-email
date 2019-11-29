package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

type Payload struct {
	EventId     string       `json:"eventID"`
	EventType   string       `json:"eventType"`
	ContentType string       `json:"contentType"`
	Data        EmailMessage `json:"data"`
}

type EmailMessage struct {
	Subject string          `json:"subject"`
	From    []*mail.Address `json:"from"`
	Message string          `json:"message"`
}

type Subscribe struct {
	Data          RequestParam `json:"data"`
	Endpoint      string       `json:"endpoint"`
	Id            string       `json:"id"`
	LastMessageId uint32
	IsTesting     bool `json:"istesting"`
}

type RequestParam struct {
	Username string `json:"username"`
	Pattern  string `json:"pattern"`
	Label    string `json:"label"`
}

var Listener = make(map[string]Subscribe)
var rtmStarted bool
var newClient *client.Client

type ReceiveHandler struct{}

func (h ReceiveHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	var password = os.Getenv("PASSWORD")
	var imapHost = os.Getenv("IMAP_HOST")
	var imapPort = os.Getenv("IMAP_PORT")

	if password == "" || imapHost == "" || imapPort == "" {
		writeErrorResponse(responseWriter, fmt.Errorf("Please provide environment variables"))
		return
	}

	decoder := json.NewDecoder(request.Body)

	var sub Subscribe
	decodeError := decoder.Decode(&sub)
	if decodeError != nil {
		writeErrorResponse(responseWriter, decodeError)
		return
	}

	log.Println("Connecting to server...")

	newClient, _ = client.DialTLS(imapHost+":"+imapPort, nil)

	log.Println("Connected")

	if err := newClient.Login(sub.Data.Username, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	res2, _ := json.Marshal(sub)
	fmt.Println(string(res2))

	Listener[sub.Data.Username] = sub
	if !rtmStarted {
		go MailRTM()
		rtmStarted = true
	}

	message := Message{"true", "Subscribed", http.StatusOK}
	writeJsonResponse(responseWriter, message)
}

func MailRTM() {
	isTest := false
	for {
		if len(Listener) > 0 {
			for k, v := range Listener {
				go getMessageUpdates(k, v)
				isTest = v.IsTesting
			}
		} else {
			rtmStarted = false
			break
		}
		time.Sleep(5 * time.Second)
		if isTest {
			break
		}
	}
}

func labelInMailbox(currentLabel string, mailboxList []string) bool {
	for _, b := range mailboxList {
		if b == currentLabel {
			return true
		}
	}
	return false
}

func getMessageUpdates(userid string, sub Subscribe) {

	selectedLabel := "INBOX"
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- newClient.List("", "*", mailboxes)
	}()

	var labelList []string
	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
		labelList = append(labelList, m.Name)
	}

	if sub.Data.Label != "" {
		check := labelInMailbox(sub.Data.Label, labelList)
		if check {
			selectedLabel = sub.Data.Label
		}
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	mBox, err := newClient.Select(selectedLabel, false)
	var receivedMessage EmailMessage

	if err != nil {
		log.Fatal(err)
	}

	if mBox.Messages == 0 {
		log.Fatal("No message in mailbox")
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(mBox.Messages)

	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := newClient.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	msg := <-messages
	if msg == nil {
		log.Fatal("Server didn't returned message")
	}

	r := msg.GetBody(&section)
	if r == nil {
		log.Fatal("Server didn't returned message body")
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Fatal(err)
	}

	header := mr.Header
	if date, err := header.Date(); err == nil {
		log.Println("Date:", date)
	}
	if from, err := header.AddressList("From"); err == nil {
		log.Println("From:", from)
		receivedMessage.From = from
	}
	if to, err := header.AddressList("To"); err == nil {
		log.Println("To:", to)
	}
	if subject, err := header.Subject(); err == nil {
		log.Println("Subject:", subject)
		receivedMessage.Subject = string(subject)
	}

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			b, _ := ioutil.ReadAll(p.Body)
			receivedMessage.Message = string(b)
		case *mail.AttachmentHeader:
			filename, _ := h.Filename()
			fmt.Println("Got attachment:", filename)
		}
	}

	datajson, datajsonerr := json.Marshal(receivedMessage)
	if datajsonerr != nil {
		log.Fatal(err)
	}
	match, matcherr := regexp.MatchString(sub.Data.Pattern, string(datajson))

	if matcherr != nil {
		log.Fatal(err)
	}

	contentType := "application/json"

	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget(sub.Endpoint),
		cloudevents.WithStructuredEncoding(),
	)
	if err != nil {
		log.Printf("failed to create transport, %v", err)
		return
	}

	c, err := cloudevents.NewClient(t,
		cloudevents.WithTimeNow(),
	)
	if err != nil {
		log.Printf("failed to create client, %v", err)
		return
	}

	source, err := url.Parse(sub.Endpoint)
	event := cloudevents.Event{
		Context: cloudevents.EventContextV01{
			EventID:     sub.Id,
			EventType:   "mail",
			Source:      cloudevents.URLRef{URL: *source},
			ContentType: &contentType,
		}.AsV01(),
		Data: receivedMessage,
	}

	if (sub.Data.Pattern == "" || match) && sub.LastMessageId != msg.SeqNum {

		sub.LastMessageId = msg.SeqNum
		Listener[userid] = sub
		resp, evt, err := c.Send(context.Background(), event)
		if err != nil {
			log.Printf("failed to send: %v (%v)", err, evt)
		}
		if resp != nil {
			fmt.Printf("Response:\n%s\n", resp)
			fmt.Printf("Got Event Response Context: %+v\n", resp)
			data := event
			if err := evt.DataAs(event); err != nil {
				fmt.Printf("Got Data Error: %s\n", err.Error())
			}
			fmt.Printf("Got Response Data: %+v\n", data)
		} else {
			log.Printf("event sent at %s", time.Now())
		}
	}
}
