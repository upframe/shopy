package pages

import (
	"strconv"

	"github.com/logpacker/PayPal-Go-SDK"

	"github.com/gorilla/sessions"
)

var c *paypalsdk.Client

// InitPayPal configures the paypal client variable
func InitPayPal(client, secret string, development bool) error {
	link := paypalsdk.APIBaseLive
	if development {
		link = paypalsdk.APIBaseSandBox
	}

	var err error
	c, err = paypalsdk.NewClient(client, secret, link)
	return err
}

var (
	// BaseAddress is the base URL of the website
	BaseAddress string
	// Templates is the path to the tempaltes folder
	Templates string
)

func displayCents(cents int) string {
	price := strconv.Itoa(cents)

	if len(price) == 1 {
		price = "0.0" + price
	} else if len(price) == 2 {
		price = "0." + price
	} else {
		cents := price[len(price)-2:]
		price = price[0:len(price)-2] + "." + cents
	}

	return price
}

var store *sessions.CookieStore

func init() {
	keyPairs := [][]byte{[]byte("HEY")}

	// Creates the new cookie session;
	store = sessions.NewCookieStore(keyPairs...)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 3,
		Secure:   false,
		HttpOnly: true,
	}
	store.Options.Domain = "localhost"
}
