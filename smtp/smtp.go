package smtp

import (
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
	auth := sasl.NewPlainClient("", email.From, s.Password)
	msg := strings.NewReader(email.BuildMessage())
	return smtp.SendMail(s.Address, auth, email.From, email.To, msg)
}
