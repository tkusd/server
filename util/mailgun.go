package util

import (
	"github.com/mailgun/mailgun-go"
	"github.com/tkusd/server/config"
)

var Mailgun mailgun.Mailgun

func init() {
	Mailgun = mailgun.NewMailgun(config.Config.Mailgun.Domain, config.Config.Mailgun.PrivateKey, config.Config.Mailgun.PublicKey)
}
