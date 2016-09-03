package email

import (
	"crypto/tls"
	"net/smtp"
)

var (
	smtpUser       string
	smtpPass       string
	smtpHost       string
	smtpPort       string
	smtpServerName string
	smtpAuth       smtp.Auth
	smtpTLSConfig  *tls.Config
	// Templates is the base path for the templates
	Templates string
	// FromDefaultEmail is the default 'From' email address
	FromDefaultEmail string
)

// InitSMTP configures the email variables
func InitSMTP(user, pass, host, port string) {
	smtpUser = user
	smtpPass = pass
	smtpHost = host
	smtpPort = port

	// Connect to the SMTP Server
	smtpServerName = smtpHost + ":" + smtpPort
	smtpAuth = smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// TLS config
	smtpTLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}
}
