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
	handler := http.HandlerFunc(SendEmail)
	handler.ServeHTTP(recorder, request)

	Describe("Send email message", func() {
		Context("send", func() {
			It("Should result http.StatusOK", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
		})
	})
})

var _ = Describe("Receive email", func() {

	receive := Received{Username: "demot636@gmail.com", Password: "Test@123"}
	requestBody := new(bytes.Buffer)
	errr := json.NewEncoder(requestBody).Encode(receive)
	if errr != nil {
		log.Fatal(errr)
	}

	request, err := http.NewRequest("POST", "/receive", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(ReceiveEmail)
	handler.ServeHTTP(recorder, request)

	Describe("receive email message", func() {
		Context("receive", func() {
			It("Should result http.StatusOK", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
		})
	})
})
