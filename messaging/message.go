package messaging

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/cloudevents/sdk-go"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	result "github.com/heaptracetechnology/microservice-email/result"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type Email struct {
	Subject string   `json:"subject,omitempty"`
	Body    string   `json:"message,omitempty"`
	From    string   `json:"from,omitempty"`
	To      []string `json:"to,omitempty"`
}

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

type Message struct {
	Success    string `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statuscode"`
}

type ErrorMessage struct {
	Success    string `json:"success"`
	Message    error  `json:"message"`
	StatusCode int    `json:"statuscode"`
}

var Listener = make(map[string]Subscribe)
var rtmStarted bool
var newClient *client.Client

//BuildMessage
func (mail *Email) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.From)
	if len(mail.To) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	}
	message += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	message += "\r\n" + mail.Body

	return message
}

func ReturnError(responseErr error) {
	var responseWriter http.ResponseWriter
	message := ErrorMessage{"false", responseErr, http.StatusBadRequest}
	bytes, _ := json.Marshal(message)
	result.WriteJsonResponse(responseWriter, bytes, http.StatusBadRequest)
	return
}

//Send Email
func Send(responseWriter http.ResponseWriter, request *http.Request) {

	var password = os.Getenv("PASSWORD")
	var smtpHost = os.Getenv("SMTP_HOST")
	var smtpPort = os.Getenv("SMTP_PORT")

	if password == "" || smtpHost == "" || smtpPort == "" {
		message := Message{"false", "Please provide environment variables", http.StatusBadRequest}
		bytes, _ := json.Marshal(message)
		result.WriteJsonResponse(responseWriter, bytes, http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(request.Body)
	var param Email
	decodeErr := decoder.Decode(&param)
	if decodeErr != nil {
		result.WriteErrorResponse(responseWriter, decodeErr)
		return
	}

	if param.From == "" || param.To == nil || param.Subject == "" || param.Body == "" {
		message := Message{"false", "Please provide required details", http.StatusBadRequest}
		bytes, _ := json.Marshal(message)
		result.WriteJsonResponse(responseWriter, bytes, http.StatusBadRequest)
		return
	}

	messageBody := param.BuildMessage()

	auth := smtp.PlainAuth("", param.From, password, smtpHost)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	smtpAddress := smtpHost + ":" + smtpPort

	conn, err := tls.Dial("tcp", smtpAddress, tlsconfig)
	if err != nil {
		ReturnError(err)
		return
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		ReturnError(err)
		return
	}

	if err = client.Auth(auth); err != nil {
		ReturnError(err)
		return
	}

	if err = client.Mail(param.From); err != nil {
		ReturnError(err)
		return
	}

	for _, k := range param.To {
		if err = client.Rcpt(k); err != nil {
			ReturnError(err)
			return
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		ReturnError(err)
		return
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		ReturnError(err)
		return
	}

	err = w.Close()
	if err != nil {
		ReturnError(err)
		return
	}

	client.Quit()

	message := Message{"true", "Mail sent successfully", http.StatusOK}
	bytes, _ := json.Marshal(message)
	result.WriteJsonResponse(responseWriter, bytes, http.StatusOK)
}

//Receiver Email
func Receiver(responseWriter http.ResponseWriter, request *http.Request) {

	var password = os.Getenv("PASSWORD")
	var imapHost = os.Getenv("IMAP_HOST")
	var imapPort = os.Getenv("IMAP_PORT")

	if password == "" || imapHost == "" || imapPort == "" {
		message := Message{"false", "Please provide environment variables", http.StatusBadRequest}
		bytes, _ := json.Marshal(message)
		result.WriteJsonResponse(responseWriter, bytes, http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(request.Body)

	var sub Subscribe
	decodeError := decoder.Decode(&sub)
	if decodeError != nil {
		result.WriteErrorResponse(responseWriter, decodeError)
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

	bytes, _ := json.Marshal("Subscribed")
	result.WriteJsonResponse(responseWriter, bytes, http.StatusOK)
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
		case mail.TextHeader:
			b, _ := ioutil.ReadAll(p.Body)
			receivedMessage.Message = string(b)
		case mail.AttachmentHeader:
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
	s1 := strings.Split(sub.Endpoint, "//")
	_, ip := s1[0], s1[1]
	s := strings.Split(ip, ":")
	_, port := s[0], s[1]
	sub.Endpoint = "http://192.168.0.61:" + string(port)

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
		resp, err := c.Send(context.Background(), event)
		if err != nil {
			log.Printf("failed to send: %v", err)
		}
		if resp != nil {
			fmt.Printf("Response:\n%s\n", resp)
			fmt.Printf("Got Event Response Context: %+v\n", resp.Context)
			data := event
			if err := resp.DataAs(event); err != nil {
				fmt.Printf("Got Data Error: %s\n", err.Error())
			}
			fmt.Printf("Got Response Data: %+v\n", data)
		} else {
			log.Printf("event sent at %s", time.Now())
		}
	}
}
