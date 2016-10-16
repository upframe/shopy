package models

import "strconv"

type CartItem struct {
	*Product
	Quantity int
}

func (i CartItem) GetPrice() int {
	return i.Price * i.Quantity
}

type Cart struct {
	Products []*CartItem
	Locked   bool
}

func (c Cart) GetTotal() int {
	total := 0

	for k := range c.Products {
		total += c.Products[k].Price * c.Products[k].Quantity
	}

	return total
}

func (c Cart) GetDescription() string {
	description := ""

	for _, product := range c.Products {
		description += strconv.Itoa(product.Quantity) + " x " + product.Name + "\n"
	}

	return description
}
