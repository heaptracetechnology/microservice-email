package http

import (
	"encoding/json"
	"net/http"

	"github.com/oms-services/email"
)

//go:generate counterfeiter . Emailer

type Emailer interface {
	Send(email email.Email) error
}

type SendHandler struct {
	Emailer Emailer

	Password string
	SMTPHost string
	SMTPPort string
}

func (h SendHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if h.Password == "" || h.SMTPHost == "" || h.SMTPPort == "" {
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

	if err := h.Emailer.Send(param); err != nil {
		message := Message{"false", err.Error(), http.StatusBadRequest}
		bytes, _ := json.Marshal(message)
		writeJsonResponse(responseWriter, bytes, http.StatusBadRequest)
		return
	}

	message := Message{"true", "Mail sent successfully", 250}
	bytes, _ := json.Marshal(message)
	writeJsonResponse(responseWriter, bytes, 250)
}
