package models

type CartItem struct {
	*Product
	Quantity int
}

func (i CartItem) GetPrice() int {
	return i.Price * i.Quantity
}

type Cart struct {
	Products []*CartItem
}

func (c Cart) GetTotal() int {
	total := 0

	for k := range c.Products {
		total += c.Products[k].Price * c.Products[k].Quantity
	}

	return total
}
