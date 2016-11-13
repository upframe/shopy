package main

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
	"github.com/upframe/fest"
	"github.com/upframe/fest/email"
	"github.com/upframe/fest/mysql"
)

type config struct {
	Development bool   `json:"Development"`
	Key         string `json:"Key"`
	InviteOnly  bool   `json:"InviteOnly"`
	BaseInvites int    `json:"BaseInvites"`
	Database    struct {
		User     string `json:"User"`
		Password string `json:"Password"`
		Host     string `json:"Host"`
		Port     string `json:"Port"`
		Name     string `json:"Name"`
	} `json:"Database"`
	SMTP struct {
		User     string `json:"User"`
		Password string `json:"Password"`
		Host     string `json:"Host"`
		Port     string `json:"Port"`
	} `json:"SMTP"`
	PayPal struct {
		Client string `json:"Client"`
		Secret string `json:"Secret"`
	} `json:"PayPal"`
}

func init() {
	// Regist types so they can be used on Cookies
	gob.Register(fest.CartCookie{})
	gob.Register(fest.OrderCookie{})

	c := &config{}

	fest.BaseInvites = c.BaseInvites
	fest.InviteOnly = c.InviteOnly

	fest.BaseAddress = "http://localhost/"
	fest.Templates = "_assets/templates"
	email.Templates = "_assets/templates/email"

	// Creates the new cookie session;
	fest.Store = sessions.NewCookieStore([]byte(c.Key))
	fest.Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 3,
		Secure:   !c.Development,
		HttpOnly: true,
	}

	fest.Store.Options.Domain = "localhost"

	// Configures the email
	email.InitSMTP(c.SMTP.User, c.SMTP.Password, c.SMTP.Host, c.SMTP.Port)

	// Connects to the database and checks for an error
	err := mysql.InitDB(
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)

	if err != nil {
		panic(err)
	}

	// Configures PayPal
	if fest.InitPayPal(c.PayPal.Client, c.PayPal.Secret, c.Development) != nil {
		panic(err)
	}
}

func main() {

}
