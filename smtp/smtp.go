package smtp

import (
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type Client struct {
	Address  string
	Password string
}

func (s Client) Send(from string, to []string, msg string) error {
	auth := sasl.NewPlainClient("", from, s.Password)
	return smtp.SendMail(s.Address, auth, from, to, strings.NewReader(msg))
}
