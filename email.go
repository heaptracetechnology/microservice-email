package email

import (
	"fmt"
	"strings"
)

type Email struct {
	Subject string   `json:"subject,omitempty"`
	Body    string   `json:"message,omitempty"`
	From    string   `json:"from,omitempty"`
	To      []string `json:"to,omitempty"`
}

func (mail *Email) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.From)
	if len(mail.To) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	}
	message += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	message += "\r\n" + mail.Body

	return message
}
