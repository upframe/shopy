package fest

import (
	"encoding/gob"
	"net/http"

	"github.com/upframe/fest"
)

// CartCookie is the cookie of the cart
type CartCookie struct {
	Products map[int]int
	Locked   bool
}

func init() {
	// Regist types so they can be used on Cookies
	gob.Register(CartCookie{})
}

func SetSessionCookie(w http.ResponseWriter, c *fest.Config, s *fest.SessionCookie) {
	if encoded, err := c.CookieStore.Encode("session", s); err == nil {
		cookie := &http.Cookie{
			Name:     "session",
			Value:    encoded,
			Path:     "/",
			MaxAge:   10800,
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}
}

func SetCartCookie(w http.ResponseWriter, c *fest.Config, s *fest.CartCookie) {
	if encoded, err := s.Encode("cart", s); err == nil {
		cookie := &http.Cookie{
			Name:     "cart",
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}
}

func ReadSessionCookie(w http.ResponseWriter, r *http.Request, c *fest.Config) {
	if cookie, err := r.Cookie("session"); err == nil {
		var value string
		if err := s.Decode("session", cookie.Value, &value); err == nil {

		}
	}
}

func ReadCartCookie(w http.ResponseWriter, r *http.Request, c *fest.Config) {
	if cookie, err := r.Cookie("cart"); err == nil {
		var value string
		if err := s.Decode("cart", cookie.Value, &value); err == nil {

		}
	}
}
