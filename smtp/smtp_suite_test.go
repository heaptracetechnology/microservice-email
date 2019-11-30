package smtp_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSmtp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Smtp Suite")
}
