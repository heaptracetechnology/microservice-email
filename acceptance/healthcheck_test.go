package acceptance_test

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Healthchecking", func() {
	var (
		session *gexec.Session
	)

	BeforeEach(func() {
		extraEnv := []string{
			fmt.Sprintf("PASSWORD=%s", getEnvOrError("EMAIL_PASSWORD")),
			fmt.Sprintf("SMTP_HOST=%s", "smtp.gmail.com"),
			fmt.Sprintf("SMTP_PORT=%s", "587"),
		}

		cmd := exec.Command(serverPath)
		cmd.Env = append(cmd.Env, extraEnv...)
		session = execBin(cmd)
	})

	AfterEach(func() {
		session.Kill().Wait()
	})

	It("eventually responds 200 OK", func() {
		Eventually(healthcheck).Should(Succeed())
	})
})
