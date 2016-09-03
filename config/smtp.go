package config

import (
	"crypto/tls"
	"net/smtp"
)

// initSMTP configures the email variables
func initSMTP() {
	// Connect to the SMTP Server
	SMTPServerName = SMTPHost + ":" + SMTPPort
	SMTPAuth = smtp.PlainAuth("", SMTPUser, SMTPPass, SMTPHost)

	// TLS config
	SMTPTLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         SMTPHost,
	}
}
