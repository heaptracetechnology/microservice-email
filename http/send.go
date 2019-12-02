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
		writeErrorResponse(responseWriter, err)
		return
	}

	writeSuccessResponse(responseWriter, "Mail sent successfully", 250)
}
