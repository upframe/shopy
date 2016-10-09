package models

import "encoding/gob"

// OrderCookie contains the information of an order
type OrderCookie struct {
	Promocode struct {
		Code           string
		DiscountAmount int
	}
	Credits int
	Total   int
}

// CartCookie is the cookie of the cart
type CartCookie struct {
	Products map[int]int
	Locked   bool
}

func init() {
	// Regist types so they can be used on Cookies
	gob.Register(CartCookie{})
	gob.Register(OrderCookie{})
}
