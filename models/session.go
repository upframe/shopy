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

	c := s.Values["Cart"].(map[int]int)

	if len(c) == 0 {
		return &Cart{}, nil
	}

	for k := range c {
		ids += strconv.Itoa(k) + ", "
	}

	ids = strings.TrimSuffix(ids, ", ") + ")"
	products := []Product{}
	err := db.Select(&products, "SELECT * FROM products WHERE id IN "+ids+" ORDER BY id")

	if err != nil {
		return nil, err
	}

	cart := &Cart{}

	for k := range products {
		cart.Products = append(cart.Products, &CartItem{
			Quantity: c[products[k].ID],
			Product:  &products[k],
		})
	}

	return cart, nil
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
