package fest

import "encoding/gob"

// CartCookie is the cookie of the cart
type CartCookie struct {
	Products map[int]int
	Locked   bool
}

func init() {
	// Regist types so they can be used on Cookies
	gob.Register(CartCookie{})
}
