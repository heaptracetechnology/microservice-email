package acceptance_test

import (
	"fmt"
	"net/http"
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

func healthcheck() error {
	resp, err := http.Get("http://localhost:3000/healthcheck")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	return fmt.Errorf("expected status code 200 but got %d", resp.StatusCode)
}
