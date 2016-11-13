package fest

import (
	"errors"

	"github.com/gorilla/sessions"
	"github.com/logpacker/PayPal-Go-SDK"
)

const (
	passwordSaltBytes = 32
	passwordHashBytes = 64
)

var (
	ErrAlreadyLoggedIn = errors.New("The user is already logged in.")
	ErrNotLoggedIn     = errors.New("The user is not logged in.")

	BaseInvites = 0
	InviteOnly  = false

	BaseAddress string
	Templates   string
	Store       *sessions.CookieStore
	PayPal      *paypalsdk.Client
)

// InitPayPal configures the paypal client variable
func InitPayPal(client, secret string, development bool) error {
	link := paypalsdk.APIBaseLive
	if development {
		link = paypalsdk.APIBaseSandBox
	}

	var err error
	PayPal, err = paypalsdk.NewClient(client, secret, link)
	return err
}
