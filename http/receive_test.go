package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/oms-services/email/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Receiving Emails", func() {

	var (
		recorder    *httptest.ResponseRecorder
		requestBody *bytes.Buffer
	)

	BeforeEach(func() {
		recorder = nil
		requestBody = &bytes.Buffer{}

		os.Unsetenv("PASSWORD")
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("IMAP_HOST")
		os.Unsetenv("IMAP_PORT")
	})

	JustBeforeEach(func() {
		request, err := http.NewRequest("POST", "/receive", requestBody)
		Expect(err).NotTo(HaveOccurred())
		recorder = httptest.NewRecorder()
		handler := ReceiveHandler{}
		handler.ServeHTTP(recorder, request)
	})

	When("all env vars are set", func() {
		BeforeEach(func() {
			os.Setenv("PASSWORD", password)
			os.Setenv("IMAP_HOST", "imap.gmail.com")
			os.Setenv("IMAP_PORT", "993")

			var received Subscribe
			var data RequestParam
			data.Username = from
			data.Pattern = "dddd"
			received.Data = data
			received.IsTesting = true

			Expect(json.NewEncoder(requestBody).Encode(received)).To(Succeed())
		})

		It("Should result http.StatusOK", func() {
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})
	})

	When("not all env vars are set", func() {
		BeforeEach(func() {
			os.Setenv("PASSWORD", password)

			var received Subscribe
			var data RequestParam
			data.Username = from
			data.Pattern = "dddd"
			received.Data = data
			received.IsTesting = true

			Expect(json.NewEncoder(requestBody).Encode(received)).To(Succeed())
		})

		It("Should result http.StatusBadRequest", func() {
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
