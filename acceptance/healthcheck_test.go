package acceptance_test

import (
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
		cmd := exec.Command(serverPath)
		session = execBin(cmd)
	})

	AfterEach(func() {
		session.Kill().Wait()
	})

	It("eventually responds 200 OK", func() {
		Eventually(healthcheck).Should(Succeed())
	})
})
