package config

import (
	"net/smtp"
)

func initSmptAuth() {
	authSmtp = smtp.PlainAuth("", smtpEmail, smtpPassword, smtpHost)
}
