package models

import "strconv"

// CartItem handles products and its quantity
type CartItem struct {
	*Product
	Quantity int
}

// GetPrice displays the product price * quantity
func (i CartItem) GetPrice() int {
	return i.Price * i.Quantity
}

// Cart does exactly what is says
type Cart struct {
	Products []*CartItem
	Locked   bool
}

// GetTotal display the order total price
func (c Cart) GetTotal() int {
	total := 0

	for k := range c.Products {
		total += c.Products[k].Price * c.Products[k].Quantity
	}

	return total
}

// GetDescription is used for the payment procedure
func (c Cart) GetDescription() string {
	description := ""

	for _, product := range c.Products {
		description += strconv.Itoa(product.Quantity) + " x " + product.Name + "\n"
	}

	return description
}
