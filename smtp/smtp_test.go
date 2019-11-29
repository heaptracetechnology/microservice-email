package smtp_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/oms-services/email/smtp"
)

var _ = Describe("SMTP Client", func() {

	var client Client

	BeforeEach(func() {
		client = Client{
			Password: getEnvOrError("EMAIL_PASSWORD"),
			Address:  "smtp.gmail.com:587",
		}
	})

	It("sends emails successfully", func() {
		from := getEnvOrError("EMAIL_FROM")
		to := getEnvOrError("EMAIL_TO")

		Expect(client.Send(from, []string{to}, "test message")).To(Succeed())
	})
})

func getEnvOrError(env string) string {
	value := os.Getenv(env)
	if value == "" {
		Fail(fmt.Sprintf("Environment variable '%s' must be set", env))
	}

	return value
}
