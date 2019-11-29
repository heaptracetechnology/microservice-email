package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/oms-services/email"
	. "github.com/oms-services/email/http"
	"github.com/oms-services/email/http/httpfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	password string
	to       string
	from     string
)

var _ = BeforeSuite(func() {
	password = getEnvOrError("EMAIL_PASSWORD")
	to = getEnvOrError("EMAIL_TO")
	from = getEnvOrError("EMAIL_FROM")
})

var _ = Describe("Sending Emails", func() {

	var (
		emailer     *httpfakes.FakeEmailer
		recorder    *httptest.ResponseRecorder
		requestBody *bytes.Buffer

		handler SendHandler
	)

	BeforeEach(func() {
		emailer = &httpfakes.FakeEmailer{}

		recorder = nil
		requestBody = &bytes.Buffer{}

		handler = SendHandler{
			Emailer: emailer,
		}

		os.Unsetenv("IMAP_HOST")
		os.Unsetenv("IMAP_PORT")
	})

	JustBeforeEach(func() {
		request, err := http.NewRequest("POST", "/send", requestBody)
		Expect(err).NotTo(HaveOccurred())

		recorder = httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
	})

	When("a valid body is sent in the request", func() {
		var emailToSend email.Email
		BeforeEach(func() {
			emailToSend = email.Email{
				From:    from,
				To:      []string{to},
				Subject: "Testing microservice",
				Body:    "Any body message to test"}

			Expect(json.NewEncoder(requestBody).Encode(emailToSend)).To(Succeed())
		})

		It("should attempt to send the email", func() {
			Expect(emailer.SendCallCount()).To(Equal(1))
			Expect(emailer.SendArgsForCall(0)).To(Equal(emailToSend))
		})

		When("emailing is successful", func() {
			BeforeEach(func() {
				emailer.SendReturns(nil)
			})

			It("should result in a successful SMTP response", func() {
				Expect(recorder.Code).To(Equal(250))
			})
		})

		When("emailing is unsuccessful", func() {
			BeforeEach(func() {
				emailer.SendReturns(errors.New("explode"))
			})

			It("should result http.StatusBadRequest", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	When("an invalid body is sent in the request", func() {
		BeforeEach(func() {
			email := []byte(`{"invalid":body}`)
			Expect(json.NewEncoder(requestBody).Encode(email)).To(Succeed())
		})

		It("should result http.StatusBadRequest", func() {
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
