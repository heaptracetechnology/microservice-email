package http

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeErrorResponse(responseWriter http.ResponseWriter, err error) {
	messageBytes, _ := json.Marshal(err)
	writeJsonResponse(responseWriter, messageBytes, http.StatusBadRequest)
}

func writeJsonResponse(responseWriter http.ResponseWriter, bytes []byte, statusCode int) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	responseWriter.WriteHeader(statusCode)
	_, err := responseWriter.Write(bytes)
	if err != nil {
		log.Println(err)
	}
}
