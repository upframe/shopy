package shopy

import (
	"net/http"
	"strconv"
	"strings"
)

// Cart does exactly what is says
type Cart struct {
	RawList  map[int]int
	Products []*CartItem
	Locked   bool
}

// CartItem handles products and its quantity
type CartItem struct {
	*Product
	Quantity int
}

// FillProducts ...
func (c *Cart) FillProducts(service ProductService) error {
	ids := "("

	if len(c.RawList) == 0 {
		return nil
	}

	for k := range c.RawList {
		ids += strconv.Itoa(k) + ", "
	}

	ids = strings.TrimSuffix(ids, ", ") + ")"
	products, err := service.GetsWhereIn(0, 0, "ID", "ID", ids)
	if err != nil {
		return err
	}

	for k := range products {
		c.Products = append(c.Products, &CartItem{
			Quantity: c.RawList[products[k].ID],
			Product:  products[k],
		})
	}

	return nil
}

// GetTotal display the order total price
func (c Cart) GetTotal() int {
	if c.Products == nil {
		return 0
	}

	total := 0

	for k := range c.Products {
		total += c.Products[k].Price * c.Products[k].Quantity
	}

	return total
}

// GetDescription is used for the payment procedure
func (c Cart) GetDescription() string {
	if c.Products == nil {
		return ""
	}

	description := ""

	for _, product := range c.Products {
		description += strconv.Itoa(product.Quantity) + " x " + product.Name + "\n"
	}

	return description
}

// GetPrice displays the product price * quantity
func (i CartItem) GetPrice() int {
	return i.Price * i.Quantity
}

// CartService ...
type CartService interface {
	Save(w http.ResponseWriter, c *Cart) error
	Get(w http.ResponseWriter, r *http.Request) (*Cart, error)
	Reset(w http.ResponseWriter) error
}
