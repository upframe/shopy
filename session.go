package fest

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
)

// Session ...
type Session struct {
	*sessions.Session
	User *User
}

// IsLoggedIn checks if the user is logged in
func (s Session) IsLoggedIn() bool {
	switch s.Values["IsLoggedIn"].(type) {
	case bool:
		return s.Values["IsLoggedIn"].(bool)
	}

	return false
}

// IsAdmin checks if an user is admin
func (s Session) IsAdmin() bool {
	if !s.IsLoggedIn() {
		return false
	}

	return s.User.Admin
}

// GetCart returns the user Cart
func (s Session) GetCart(service ProductService) (*Cart, error) {
	ids := "("

	c := s.Values["Cart"].(CartCookie)

	if len(c.Products) == 0 {
		return &Cart{}, nil
	}

	for k := range c.Products {
		ids += strconv.Itoa(k) + ", "
	}

	ids = strings.TrimSuffix(ids, ", ") + ")"
	products, err := service.GetsWhereIn(0, 0, "ID", "ID", ids)
	if err != nil {
		return nil, err
	}

	cart := &Cart{}
	cart.Locked = c.Locked

	for k := range products {
		cart.Products = append(cart.Products, &CartItem{
			Quantity: c.Products[products[k].ID],
			Product:  products[k],
		})
	}

	return cart, nil
}

// SessionService ...
type SessionService interface {
	Session(w http.ResponseWriter, r *http.Request) (*Session, error)
}
