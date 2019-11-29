package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/oms-services/email/smtp"
)

type Email struct {
	Subject string   `json:"subject,omitempty"`
	Body    string   `json:"message,omitempty"`
	From    string   `json:"from,omitempty"`
	To      []string `json:"to,omitempty"`
}

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

type SendHandler struct{}

func (h SendHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	var password = os.Getenv("PASSWORD")
	var smtpHost = os.Getenv("SMTP_HOST")
	var smtpPort = os.Getenv("SMTP_PORT")

	if password == "" || smtpHost == "" || smtpPort == "" {
		message := Message{"false", "Please provide environment variables", http.StatusBadRequest}
		bytes, _ := json.Marshal(message)
		writeJsonResponse(responseWriter, bytes, http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(request.Body)
	var param Email
	decodeErr := decoder.Decode(&param)
	if decodeErr != nil {
		writeErrorResponse(responseWriter, decodeErr)
		return
	}

	if param.From == "" || param.To == nil || param.Subject == "" || param.Body == "" {
		message := Message{"false", "Please provide required details", http.StatusBadRequest}
		bytes, _ := json.Marshal(message)
		writeJsonResponse(responseWriter, bytes, http.StatusBadRequest)
		return
	}

	messageBody := param.BuildMessage()

	smtpAddress := smtpHost + ":" + smtpPort
	client := smtp.Client{
		Address:  smtpAddress,
		Password: password,
	}

	if err := client.Send(param.From, param.To, messageBody); err != nil {
		message := Message{"false", err.Error(), http.StatusBadRequest}
		bytes, _ := json.Marshal(message)
		writeJsonResponse(responseWriter, bytes, http.StatusBadRequest)
		return
	}

	message := Message{"true", "Mail sent successfully", 250}
	bytes, _ := json.Marshal(message)
	writeJsonResponse(responseWriter, bytes, 250)
}
