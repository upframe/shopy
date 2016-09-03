package models

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"path/filepath"
)

// FromDefaultEmail is the default FROM email
const FromDefaultEmail = "noreply@bitsn.me"

var (
	smtpUser       string
	smtpPass       string
	smtpHost       string
	smtpPort       string
	smtpServerName string
	smtpAuth       smtp.Auth
	smtpTLSConfig  *tls.Config

	// TemplatesPath is where are the emails templates placed
	TemplatesPath string
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

// Email contains the information about email
type Email struct {
	From    *mail.Address
	To      *mail.Address
	Subject string
	Body    string
}

// UseTemplate adds  the tempalte to the email and renders it to the Body field
func (e *Email) UseTemplate(name string, data interface{}) error {
	// Opens the template file and checks if there is any error
	page, err := ioutil.ReadFile(filepath.Clean(TemplatesPath + name + ".tmpl"))
	if err != nil {
		return err
	}

	funcs := template.FuncMap{
		"CSS": func(s string) template.CSS {
			return template.CSS(s)
		},
		"HTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"URL": func(s string) template.URL {
			return template.URL(s)
		},
	}

	// Create the template with the content of the file
	tpl, err := template.New(name).Funcs(funcs).Parse(string(page))

	if err != nil {
		return err
	}

	// Creates a buffer to render the template into it
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, data)

	if err != nil {
		return err
	}

	e.Body = buf.String()
	return nil
}

// Send sends the email
func (e Email) Send() error {
	// Setup headers
	headers := make(map[string]string)
	headers["From"] = e.From.String()
	headers["To"] = e.To.String()
	headers["Subject"] = e.Subject
	headers["Content-Type"] = "text/html"

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + e.Body

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", smtpServerName, smtpTLSConfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(smtpAuth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(e.From.Address); err != nil {
		return err
	}

	if err = c.Rcpt(e.To.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()
	return nil
}
