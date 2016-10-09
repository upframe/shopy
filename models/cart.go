package models

type CartItem struct {
	*Product
	Quantity int
}

type Cart struct {
	Products []*CartItem
}

func (c Cart) GetTotal() int {
	return 0
}
