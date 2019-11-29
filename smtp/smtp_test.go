package smtp_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/oms-services/email"
	. "github.com/oms-services/email/smtp"
)

var _ = Describe("SMTP Client", func() {

	var (
		client Client

		emailToSend email.Email
		err         error
	)

	BeforeEach(func() {
		client = Client{
			Password: getEnvOrError("EMAIL_PASSWORD"),
			Address:  "smtp.gmail.com:587",
		}

		from := getEnvOrError("EMAIL_FROM")
		to := getEnvOrError("EMAIL_TO")

		emailToSend = email.Email{
			Subject: "Test Subject",
			Body:    "Test Body",
			From:    from,
			To:      []string{to},
		}
	})

	JustBeforeEach(func() {
		err = client.Send(emailToSend)
	})

	PWhen("a valid email is provided", func() {
		It("sends emails successfully", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("the email does not contain a subject", func() {
		BeforeEach(func() {
			emailToSend.Subject = ""
		})

		It("returns an error", func() {
			Expect(err).To(MatchError("please provide a subject for the email"))
		})
	})

	When("the email does not contain a body", func() {
		BeforeEach(func() {
			emailToSend.Body = ""
		})

		It("returns an error", func() {
			Expect(err).To(MatchError("please provide a body for the email"))
		})
	})

	When("the email does not contain a from email", func() {
		BeforeEach(func() {
			emailToSend.From = ""
		})

		It("returns an error", func() {
			Expect(err).To(MatchError("please provide an email address to send the email from"))
		})
	})

	When("the email does not contain a to email", func() {
		BeforeEach(func() {
			emailToSend.To = nil
		})

		It("returns an error", func() {
			Expect(err).To(MatchError("please provide an email address to send the email to"))
		})
	})
})

func getEnvOrError(env string) string {
	value := os.Getenv(env)
	if value == "" {
		Fail(fmt.Sprintf("Environment variable '%s' must be set", env))
	}

	return value
}
