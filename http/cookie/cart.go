package cookie

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/bruhs/shopy"
)

// SessionCookie ...
type cart struct {
	Products map[int]int
	Locked   bool
}

func init() {
	gob.Register(&cart{})
}

// CartService ...
type CartService struct {
	Store  *securecookie.SecureCookie
	Secure bool
}

// Save ...
func (s *CartService) Save(w http.ResponseWriter, c *shopy.Cart) error {
	encoded, err := s.Store.Encode("cart", &cart{Products: c.RawList, Locked: c.Locked})
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "cart",
		Value:    encoded,
		Path:     "/",
		Secure:   s.Secure,
		HttpOnly: true,
		MaxAge:   60 * 60 * 24 * 365,
	}

	http.SetCookie(w, cookie)
	return nil
}

// Get ...
func (s *CartService) Get(w http.ResponseWriter, r *http.Request) (*shopy.Cart, error) {
	c := &shopy.Cart{Locked: false, RawList: map[int]int{}}

	cookie, err := r.Cookie("cart")
	if err != nil {
		return c, s.Reset(w)
	}

	var value *cart
	// if the decryption keys aren't right
	err = s.Store.Decode("cart", cookie.Value, &value)
	if err != nil {
		return c, s.Reset(w)
	}

	c.Locked = value.Locked
	c.RawList = value.Products
	return c, nil
}

// Reset ...
func (s *CartService) Reset(w http.ResponseWriter) error {
	return s.Save(w, &shopy.Cart{Locked: false, RawList: map[int]int{}})
}
