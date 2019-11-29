package smtp

import (
	"fmt"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/oms-services/email"
)

type Client struct {
	Address  string
	Password string
}

func (s Client) Send(email email.Email) error {
	if err := validate(email); err != nil {
		return err
	}

	auth := sasl.NewPlainClient("", email.From, s.Password)
	msg := strings.NewReader(email.BuildMessage())
	return smtp.SendMail(s.Address, auth, email.From, email.To, msg)
}

func validate(email email.Email) error {
	if email.Subject == "" {
		return fmt.Errorf("please provide a subject for the email")
	}

	if email.Body == "" {
		return fmt.Errorf("please provide a body for the email")
	}

	if email.From == "" {
		return fmt.Errorf("please provide an email address to send the email from")
	}

	if len(email.To) == 0 {
		return fmt.Errorf("please provide an email address to send the email to")
	}

	return nil
}
