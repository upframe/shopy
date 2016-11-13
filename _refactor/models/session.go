package models

import (
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
)

// Session wraps the gorilla session and adds some more useful functions and
// information in the context of the application such as the user information
// and cart of products
type Session struct {
	*sessions.Session
	User *User
}

// GetCart returns the user Cart
func (s Session) GetCart() (*Cart, error) {
	ids := "("

	c := s.Values["Cart"].(CartCookie)

	if len(c.Products) == 0 {
		return &Cart{}, nil
	}

	for k := range c.Products {
		ids += strconv.Itoa(k) + ", "
	}

	ids = strings.TrimSuffix(ids, ", ") + ")"
	products := []Product{}
	err := db.Select(&products, "SELECT * FROM products WHERE id IN "+ids+" ORDER BY id")

	if err != nil {
		return nil, err
	}

	cart := &Cart{}
	cart.Locked = c.Locked

	for k := range products {
		cart.Products = append(cart.Products, &CartItem{
			Quantity: c.Products[products[k].ID],
			Product:  &products[k],
		})
	}

	return cart, nil
}
