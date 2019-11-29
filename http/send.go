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
}

func (h SendHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var param email.Email
	if err := decoder.Decode(&param); err != nil {
		writeErrorResponse(responseWriter, err)
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
