package messaging

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	password = os.Getenv("EMAIL_PASSWORD")
	to       = os.Getenv("EMAIL_TO")
	from     = os.Getenv("EMAIL_FROM")
)

//Negative test without enviroment variables send mail
var _ = Describe("Send email", func() {

	to := []string{to}
	email := Email{From: from, To: to, Subject: "Testing microservice", Body: "Any body message to test"}
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
				Expect(http.StatusBadRequest).To(Equal(http.StatusBadRequest))
			})
		})
	})
})

//Negative test without from variables send mail
var _ = Describe("Send email", func() {

	os.Setenv("PASSWORD", password)
	os.Setenv("SMTP_HOST", "smtp.gmail.com")
	os.Setenv("SMTP_PORT", "465")

	to := []string{to}
	email := Email{From: from, To: to, Subject: "Testing microservice", Body: "Any body message to test"}
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
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})

//Negative test without smtp variables send mail
var _ = Describe("Send email", func() {

	os.Setenv("PASSWORD", password)
	os.Setenv("SMTP_PORT", "465")

	to := []string{to}
	email := Email{From: from, To: to, Subject: "Testing microservice", Body: "Any body message to test"}
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
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})

//Negative test without args variables send mail
var _ = Describe("Send email", func() {

	os.Setenv("PASSWORD", password)
	os.Setenv("SMTP_HOST", "smtp.gmail.com")
	os.Setenv("SMTP_PORT", "465")

	email := Email{}
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
				Expect(http.StatusBadRequest).To(Equal(http.StatusBadRequest))
			})
		})
	})
})

//Postive send mail test
var _ = Describe("Send email", func() {

	os.Setenv("PASSWORD", password)
	os.Setenv("SMTP_HOST", "smtp.gmail.com")
	os.Setenv("SMTP_PORT", "465")

	to := []string{to}
	email := Email{From: from, To: to, Subject: "Testing microservice", Body: "Any body message to test"}
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
				Expect(recorder.Code).To(Equal(250))
			})
		})
	})
})

//Decoder test
var _ = Describe("Send email", func() {

	os.Setenv("PASSWORD", password)
	os.Setenv("SMTP_HOST", "smtp.gmail.com")
	os.Setenv("SMTP_PORT", "465")

	email := []byte(`{"status":false}`)
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
				Expect(http.StatusBadRequest).To(Equal(http.StatusBadRequest))
			})
		})
	})
})

//Received email negative
var _ = Describe("Received email", func() {

	os.Setenv("PASSWORD", password)

	var received Subscribe
	var data RequestParam
	data.Username = from
	data.Pattern = "dddd"
	received.Data = data
	received.IsTesting = true

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
			It("Should result http.StatusBadRequest", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})

//Received email
var _ = Describe("Received email", func() {

	os.Setenv("PASSWORD", password)
	os.Setenv("IMAP_HOST", "imap.gmail.com")
	os.Setenv("IMAP_PORT", "993")

	var received Subscribe
	var data RequestParam
	data.Username = from
	data.Pattern = "dddd"
	received.Data = data
	received.IsTesting = true

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
