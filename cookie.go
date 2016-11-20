package fest

import (
	"encoding/gob"
	"strconv"
	"strings"
)

// CartCookie is the cookie of the cart
type CartCookie struct {
	Products map[int]int
	Locked   bool
}

// GetCart returns the user Cart
func (c *CartCookie) GetCart(service ProductService) (*Cart, error) {
	ids := "("

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

// SessionCookie ...
type SessionCookie struct {
	Logged bool
	UserID int
	user   *User
}

// User ...
func (s *SessionCookie) User() *User {
	return s.user
}

// SetUser ...
func (s *SessionCookie) SetUser(u *User) {
	s.user = u
}

func init() {
	// Regist types so they can be used on Cookies
	gob.Register(CartCookie{})
	gob.Register(SessionCookie{})
}
