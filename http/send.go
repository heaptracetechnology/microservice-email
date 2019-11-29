package http

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/oms-services/email"
	"github.com/oms-services/email/smtp"
)

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
	var param email.Email
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
