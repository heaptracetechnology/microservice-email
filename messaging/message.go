package messaging

import (
	b "bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	result "github.com/heaptracetechnology/microservice-mail/result"
)

type Email struct {
	Subject  string `json:"subject,omitempty"`
	Body     string `json:"message,omitempty"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Password string `json:"password,omitempty"`
	SMTPHost string `json:"smtp_host,omitempty"`
	SMTPPort string `json:"smtp_port,omitempty"`
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
	Password string `json:"password"`
	Pattern  string `json:"pattern"`
	Label    string `json:"label"`
	ImapHost string `json:"imap_host"`
	ImapPort string `json:"imap_port"`
}

type Message struct {
	Success    string `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statuscode"`
}

var Listener = make(map[string]Subscribe)
var rtmStarted bool
var newClient *client.Client

//Send Email
func Send(responseWriter http.ResponseWriter, request *http.Request) {

	responseWriter.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(request.Body)
	var param Email
	decodeErr := decoder.Decode(&param)
	if decodeErr != nil {
		result.WriteErrorResponse(responseWriter, decodeErr)
		return
	}

	from := param.From
	to := param.To
	sub := param.Subject
	body := param.Body
	smtpAddress := param.SMTPHost + ":" + param.SMTPPort

	msg := "From: " + from + "\n" + "To: " + to + "\n" + "Subject: " + sub + "\n" + body

	err := smtp.SendMail(smtpAddress, smtp.PlainAuth("", from, param.Password, param.SMTPHost), from, []string{to}, []byte(msg))
	if err != nil {
		fmt.Println("err ::", err)
		return
	} else {
		message := Message{"true", "Email sent", http.StatusOK}
		bytes, _ := json.Marshal(message)
		result.WriteJsonResponse(responseWriter, bytes, http.StatusOK)
	}

}

//Receiver Email
func Receiver(responseWriter http.ResponseWriter, request *http.Request) {

	responseWriter.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(request.Body)

	var sub Subscribe
	decodeError := decoder.Decode(&sub)
	if decodeError != nil {
		result.WriteErrorResponse(responseWriter, decodeError)
		return
	}

	log.Println("Connecting to server...")

	newClient, _ = client.DialTLS(sub.Data.ImapHost+":"+sub.Data.ImapPort, nil)

	log.Println("Connected")

	if err := newClient.Login(sub.Data.Username, sub.Data.Password); err != nil {
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
	var data EmailMessage

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
		data.From = from
	}
	if to, err := header.AddressList("To"); err == nil {
		log.Println("To:", to)
	}
	if subject, err := header.Subject(); err == nil {
		log.Println("Subject:", subject)
		data.Subject = string(subject)
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
			data.Message = string(b)
		case mail.AttachmentHeader:
			filename, _ := h.Filename()
			fmt.Println("Got attachment:", filename)
		}
	}

	datajson, datajsonerr := json.Marshal(data)
	if datajsonerr != nil {
		log.Fatal(err)
	}
	match, matcherr := regexp.MatchString(sub.Data.Pattern, string(datajson))

	if matcherr != nil {
		log.Fatal(err)
	}

	if (sub.Data.Pattern == "" || match) && sub.LastMessageId != msg.SeqNum {
		sub.LastMessageId = msg.SeqNum

		Listener[userid] = sub
		var payload Payload
		payload.ContentType = "application" + "/" + "json"
		payload.EventType = "mail"
		payload.EventId = sub.Id
		payload.Data = data

		requestBody := new(b.Buffer)
		encodeError := json.NewEncoder(requestBody).Encode(payload)
		if encodeError != nil {
			log.Fatalln(encodeError)
			fmt.Println("err :", encodeError)
		}
		s1 := strings.Split(sub.Endpoint, "//")
		_, ip := s1[0], s1[1]
		s := strings.Split(ip, ":")
		_, port := s[0], s[1]
		sub.Endpoint = "http://192.168.1.61:" + string(port)
		res, reserror := http.Post(sub.Endpoint, "application/json", requestBody)
		if reserror != nil {
			fmt.Println("err :", reserror)
		}
		fmt.Println("res :", res)
	}
}
