package messaging

import (
	"bytes"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

var _ = Describe("Send email", func() {

	os.Setenv("PASSWORD", "ltihivyeggcimelm")
	os.Setenv("SMTP_HOST", "smtp.gmail.com")
	os.Setenv("SMTP_PORT", "465")

	to := []string{"rohits@heaptrace.com"}
	email := Email{From: "demot636@gmail.com", To: to, Subject: "Testing microservice", Body: "Any body message to test"}
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

// var _ = Describe("Received email", func() {

// 	os.Setenv("PASSWORD", "ltihivyeggcimelm")
// 	os.Setenv("IMAP_HOST", "imap.gmail.com")
// 	os.Setenv("IMAP_PORT", "993")

// 	var received Subscribe
// 	var data RequestParam
// 	data.Username = "demot636@gmail.com"
// 	data.Pattern = "dddd"
// 	received.Data = data

// 	requestBody := new(bytes.Buffer)
// 	errr := json.NewEncoder(requestBody).Encode(received)
// 	if errr != nil {
// 		log.Fatal(errr)
// 	}

// 	request, err := http.NewRequest("POST", "/receive", requestBody)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	recorder := httptest.NewRecorder()
// 	handler := http.HandlerFunc(Receiver)
// 	handler.ServeHTTP(recorder, request)

// 	Describe("received email message", func() {
// 		Context("received", func() {
// 			It("Should result http.StatusOK", func() {
// 				Expect(recorder.Code).To(Equal(http.StatusOK))
// 			})
// 		})
// 	})
// })
