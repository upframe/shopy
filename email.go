package shopy

import "net/mail"

// Email contains the information about email
type Email struct {
	From    *mail.Address
	To      *mail.Address
	Subject string
	Body    string
}

// EmailService ...
type EmailService interface {
	UseTemplate(e *Email, data interface{}, template string) error
	Send(e *Email) error
}
