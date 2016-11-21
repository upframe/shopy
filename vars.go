package fest

import (
	"errors"
	"log"

	"github.com/gorilla/securecookie"
	"github.com/logpacker/PayPal-Go-SDK"
)

// UpdateAll is used as a placeholder to update all of the fields
const UpdateAll = "#update#"

const (
	passwordSaltBytes = 32
	passwordHashBytes = 64
)

var (
	// ErrAlreadyLoggedIn ...
	ErrAlreadyLoggedIn = errors.New("The user is already logged in.")
	// ErrNotLoggedIn ...
	ErrNotLoggedIn = errors.New("The user is not logged in.")
	// ErrNotFound ...
	ErrNotFound = errors.New("Not found.")
)

// Config ...
type Config struct {
	DefaultInvites int
	InviteOnly     bool
	Port           string
	BaseAddress    string
	Templates      string
	Domain         string
	Scheme         string
	Assets         string
	CookieStore    *securecookie.SecureCookie
	PayPal         *paypalsdk.Client
	Logger         *log.Logger
	Services       *Services
}

// Services ...
type Services struct {
	Order     OrderService
	Product   ProductService
	Promocode PromocodeService
	User      UserService
	Link      LinkService
	Email     EmailService
}

// InitPayPal configures the paypal client variable
func InitPayPal(client, secret string, development bool) (*paypalsdk.Client, error) {
	link := paypalsdk.APIBaseLive
	if development {
		link = paypalsdk.APIBaseSandBox
	}

	paypal, err := paypalsdk.NewClient(client, secret, link)
	return paypal, err
}
