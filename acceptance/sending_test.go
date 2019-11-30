package acceptance_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/oms-services/email"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Sending Emails", func() {

	var session *gexec.Session

	BeforeEach(func() {
		extraEnv := []string{
			fmt.Sprintf("PASSWORD=%s", getEnvOrError("EMAIL_PASSWORD")),
			fmt.Sprintf("SMTP_HOST=%s", "smtp.gmail.com"),
			fmt.Sprintf("SMTP_PORT=%s", "587"),
		}

		cmd := exec.Command(serverPath)
		cmd.Env = append(cmd.Env, extraEnv...)
		session = execBin(cmd)

		Eventually(healthcheck).Should(Succeed())
	})

	AfterEach(func() {
		session.Kill().Wait()
	})

	When("a valid body is sent in the request", func() {

		It("should result in a successful SMTP response", func() {
			to := getEnvOrError("EMAIL_TO")
			from := getEnvOrError("EMAIL_FROM")

			email := email.Email{
				From:    from,
				To:      []string{to},
				Subject: "Testing microservice",
				Body:    "Any body message to test"}

			requestBody := &bytes.Buffer{}
			Expect(json.NewEncoder(requestBody).Encode(email)).To(Succeed())

			resp, err := http.Post("http://localhost:3000/send", "application/json", requestBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Body.Close()).To(Succeed())

			Expect(resp.StatusCode).To(Equal(250))
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
