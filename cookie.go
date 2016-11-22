package fest

import (
	"encoding/gob"
	"net/http"
	"strconv"
	"strings"
)

// CartCookie is the cookie of the cart
type CartCookie struct {
	Products map[int]int
	Locked   bool
}

// GetCart returns the user Cart
func (c *CartCookie) GetCart(service ProductService) (*Cart, error) {
	ids := "("

	if len(c.Products) == 0 {
		return &Cart{}, nil
	}

	for k := range c.Products {
		ids += strconv.Itoa(k) + ", "
	}

	ids = strings.TrimSuffix(ids, ", ") + ")"
	products, err := service.GetsWhereIn(0, 0, "ID", "ID", ids)
	if err != nil {
		return nil, err
	}

	cart := &Cart{}
	cart.Locked = c.Locked

	for k := range products {
		cart.Products = append(cart.Products, &CartItem{
			Quantity: c.Products[products[k].ID],
			Product:  products[k],
		})
	}

	return cart, nil
}

// Session ...
type Session struct {
	Logged bool
	User   *User
}

// SessionService ...
type SessionService interface {
	Save(w http.ResponseWriter, sess *Session) error
	Get(w http.ResponseWriter, r *http.Request) (*Session, error)
	Reset(w http.ResponseWriter) error
}

func init() {
	// Regist types so they can be used on Cookies
	gob.Register(CartCookie{})
}
