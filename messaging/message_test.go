package messaging

import (
	"bytes"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Send email", func() {

	email := Email{From: "demot636@gmail.com", Password: "Test@123", To: "rohits@heaptrace.com", Subject: "Testing microservice", Body: "Any body message to test", SMTPHost: "smtp.gmail.com", SMTPPort: "587"}
	requestBody := new(bytes.Buffer)
	errr := json.NewEncoder(requestBody).Encode(email)
	if errr != nil {
		log.Fatal(errr)
	}

	request, err := http.NewRequest("POST", "/send", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Send)
	handler.ServeHTTP(recorder, request)

	Describe("Send email message", func() {
		Context("send", func() {
			It("Should result http.StatusOK", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
		})
	})
})

var _ = Describe("Received email", func() {

	var received Subscribe
	var data RequestParam
	data.Username = "demot636@gmail.com"
	data.Password = "Test@123"
	data.Pattern = "dddd"
	data.ImapHost = "imap.gmail.com"
	data.ImapPort = "993"
	received.Data = data

	requestBody := new(bytes.Buffer)
	errr := json.NewEncoder(requestBody).Encode(received)
	if errr != nil {
		log.Fatal(errr)
	}

	request, err := http.NewRequest("POST", "/receive", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Receiver)
	handler.ServeHTTP(recorder, request)

	Describe("received email message", func() {
		Context("received", func() {
			It("Should result http.StatusOK", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
		})
	})
})
