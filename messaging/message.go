package messaging

import (
	"encoding/json"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	result "github.com/heaptracetechnology/microservice-mail/result"
	//"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
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

type Subscribe struct {
	Data     Data   `json:"data"`
	Endpoint string `json:"endpoint"`
	Id       string `json:"id"`
}

type Data struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Received struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	Success    string `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statuscode"`
}

//Send Email
func Send(responseWriter http.ResponseWriter, request *http.Request) {

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

//Receive Email
func Receive(responseWriter http.ResponseWriter, request *http.Request) {

	// req, _ := json.Marshal(request.Body)
	// fmt.Println("req :: ", req)

	// var param Received

	// body, err := ioutil.ReadAll(request.Body)
	// if err != nil {
	// 	fmt.Println("err ::", err)
	// }
	// defer request.Body.Close()
	// er := json.Unmarshal(body, &param)
	// if er != nil {
	// 	fmt.Println("er ::", er)
	// }

	decoder := json.NewDecoder(request.Body)
	var param Received
	decodeErr := decoder.Decode(&param)
	if decodeErr != nil {
		result.WriteErrorResponse(responseWriter, decodeErr)
		return
	}
	log.Println("Connecting to server...")

	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	defer c.Logout()

	username := "demot636@gmail.com"
	password := "Test@123"

	if err := c.Login(username, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 3 {
		from = mbox.Messages - 3
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	log.Println("Last 4 messages:")
	for msg := range messages {
		log.Println("Email Subject: " + msg.Envelope.Subject)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
}

func getEmailUpdates() {

}
