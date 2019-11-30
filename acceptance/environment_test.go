package acceptance_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Environment Prerequisities", func() {
	var (
		extraEnv []string

		session *gexec.Session
	)

	JustBeforeEach(func() {
		cmd := exec.Command(serverPath)
		cmd.Env = append(cmd.Env, extraEnv...)
		session = execBin(cmd)
	})

	AfterEach(func() {
		session.Kill().Wait()
	})

	When("no SMTP Host is provided", func() {
		BeforeEach(func() {
			extraEnv = []string{
				"PASSWORD=FOO",
				"SMTP_PORT=BAZ",
			}
		})

		It("prints an informative message to stderr", func() {
			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say(`Environment variable 'SMTP_HOST' must be set`))
		})
	})

	When("no SMTP Port is provided", func() {
		BeforeEach(func() {
			extraEnv = []string{
				"PASSWORD=FOO",
				"SMTP_HOST=BAR",
			}
		})

		It("prints an informative message to stderr", func() {
			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say(`Environment variable 'SMTP_PORT' must be set`))
		})
	})

	When("no Password is provided", func() {
		BeforeEach(func() {
			extraEnv = []string{
				"SMTP_HOST=BAR",
				"SMTP_PORT=BAZ",
			}
		})

		It("prints an informative message to stderr", func() {
			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say(`Environment variable 'PASSWORD' must be set`))
		})
	})
})
