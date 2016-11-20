package fest

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/logpacker/PayPal-Go-SDK"
)

type config struct {
	Development    bool
	Key1           string
	Key2           string
	Domain         string
	Port           int
	Scheme         string
	Assets         string
	InviteOnly     bool
	DefaultInvites int
	Database       struct {
		User     string
		Password string
		Host     string
		Port     string
		Name     string
	}
	SMTP struct {
		User     string
		Password string
		Host     string
		Port     string
	}
	PayPal struct {
		Client string
		Secret string
	}
}

// ConfigFile ...
func ConfigFile(path string) (*config, error) {
	file := &config{}

	configFile, err := os.Open("config.json")
	if err != nil {
		return file, err
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&file)
	return file, nil
}

// UpdateAll is used as a placeholder to update all of the fields
const UpdateAll = "#update#"

const (
	passwordSaltBytes = 32
	passwordHashBytes = 64
)

var (
	ErrAlreadyLoggedIn = errors.New("The user is already logged in.")
	ErrNotLoggedIn     = errors.New("The user is not logged in.")
	ErrNotFound        = errors.New("Not found.")
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
