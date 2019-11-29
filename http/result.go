package http

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Success    string `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func writeErrorResponse(responseWriter http.ResponseWriter, err error) {
	message := Message{
		Success:    "false",
		Message:    err.Error(),
		StatusCode: http.StatusBadRequest,
	}

	writeJsonResponse(responseWriter, message)
}

func writeJsonResponse(responseWriter http.ResponseWriter, message Message) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	responseWriter.WriteHeader(message.StatusCode)
	bytes, _ := json.Marshal(message)
	responseWriter.Write(bytes)
}
