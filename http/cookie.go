package http

import (
	"net/http"

	"github.com/upframe/fest"
)

// SetCartCookie ...
func SetCartCookie(w http.ResponseWriter, c *fest.Config, cart *fest.CartCookie) error {
	encoded, err := c.CookieStore.Encode("cart", cart)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "cart",
		Value:    encoded,
		Path:     "/",
		Secure:   c.Scheme == "https",
		HttpOnly: true,
		MaxAge:   60 * 60 * 24 * 365,
	}

	http.SetCookie(w, cookie)
	return nil
}

// ReadCartCookie ...
func ReadCartCookie(w http.ResponseWriter, r *http.Request, c *fest.Config) (*fest.CartCookie, error) {
	reset := func() (*fest.CartCookie, error) {
		s := &fest.CartCookie{Products: map[int]int{}}
		err := SetCartCookie(w, c, s)
		return s, err
	}

	cookie, err := r.Cookie("cart")
	if err != nil {
		return reset()
	}

	var value *fest.CartCookie
	// if the decryption keys aren't right
	err = c.CookieStore.Decode("cart", cookie.Value, &value)
	if err != nil {
		return reset()
	}

	return value, nil
}
