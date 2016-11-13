package fest

// OrderCookie contains the information of an order
type OrderCookie struct {
	Promocode struct {
		Code           string
		DiscountAmount int
		ID             int
	}
	Credits int
	Total   int
}

// CartCookie is the cookie of the cart
type CartCookie struct {
	Products map[int]int
	Locked   bool
}
