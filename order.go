package fest

// Order ...
type Order struct {
	ID        int
	UserID    int
	PayPal    string
	Value     int
	Status    string
	Credits   int
	Promocode *Promocode
	Products  []*OrderProduct
}

// OrderProduct ...
type OrderProduct struct {
	ID       int
	Name     string
	Quantity int `db:"quantity"`
}

// OrderService ...
type OrderService interface {
	Get(id int) (*Order, error)
	Gets(first, limit int, order string) ([]*Order, error)
	GetsWhere(first, limit int, order, where, sth string) ([]*Order, error)

	Create(o *Order) error
	Update(o *Order, fields ...string) error
	Delete(id int) error
}
