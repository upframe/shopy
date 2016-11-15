package fest

import (
	"errors"
	"log"

	"github.com/gorilla/sessions"
	"github.com/logpacker/PayPal-Go-SDK"
)

// UpdateAll is used as a placeholder to update all of the fields
const UpdateAll = "#update#"

const (
	passwordSaltBytes = 32
	passwordHashBytes = 64
)

var (
	ErrAlreadyLoggedIn = errors.New("The user is already logged in.")
	ErrNotLoggedIn     = errors.New("The user is not logged in.")
)

// Config ...
type Config struct {
	DefaultInvites int
	InviteOnly     bool
	BaseAddress    string
	Templates      string
	Store          *sessions.CookieStore
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
