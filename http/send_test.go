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

var _ = FDescribe("Sending Emails", func() {

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

			Password: password,
			SMTPHost: "smtp.gmail.com",
			SMTPPort: "587",
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

	When("all configuration is set correctly", func() {
		When("a valid body is sent in the request", func() {
			BeforeEach(func() {
				email := email.Email{
					From:    from,
					To:      []string{to},
					Subject: "Testing microservice",
					Body:    "Any body message to test"}

				Expect(json.NewEncoder(requestBody).Encode(email)).To(Succeed())
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

		When("the body does not contain required details", func() {
			BeforeEach(func() {
				email := email.Email{}
				Expect(json.NewEncoder(requestBody).Encode(email)).To(Succeed())
			})

			It("should result http.StatusBadRequest", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	When("not all configuration is set correctly", func() {
		When("no configuration is set", func() {
			BeforeEach(func() {
				handler.Password = ""
				handler.SMTPHost = ""
				handler.SMTPPort = ""
			})

			It("Should result http.StatusOK", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("no smtp host is set", func() {
			BeforeEach(func() {
				handler.SMTPHost = ""
			})

			It("Should result http.StatusBadRequest", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
