package config

import (
	"crypto/tls"
	"net/smtp"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

var (
	// db is our database connection
	db *sqlx.DB
	// Store stores the session cookies and help us to handle them
	Store *sessions.CookieStore
	// BaseAddress is the base URL to build URLs
	BaseAddress string
	// RootPath is the 'root' directive defined in Caddyfile
	RootPath string
	// TemplatesPath is where the templates are stored
	TemplatesPath string
	// SMTPUser is the user to connect to the SMTP server
	SMTPUser string
	// SMTPPass is the pass to connect to the SMTP server
	SMTPPass string
	// SMTPHost is the host of the SMTP server
	SMTPHost string
	// SMTPPort is port to connect to the SMTP server
	SMTPPort string
	// SMTPServerName is the name of the SMTP server
	SMTPServerName string
	// SMTPAuth contains the authentication to the SMTP server
	SMTPAuth smtp.Auth
	// SMTPTLSConfig is the TLS configuration to the SMTP server
	SMTPTLSConfig *tls.Config
	// FromDefaultEmail is the default email to send messages
	FromDefaultEmail = "noreply@bitsn.me"
	// Database connection variables
	dbUser string
	dbPass string
	dbHost string
	dbName string
	dbPort = "3306"
)
