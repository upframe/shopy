package fest

// Order ...
type Order struct {
	ID          int             `db:"id"`
	UserID      int             `db:"user_id"`
	PromocodeID NullInt64       `db:"promocode_id"`
	PayPalID    string          `db:"paypal_id"`
	Value       int             `db:"value"`
	Status      string          `db:"status"`
	Credits     int             `db:"credits"`
	Promocode   *Promocode      `db:"-"`
	Products    []*OrderProduct `db:"-"`
}

// OrderProduct ...
type OrderProduct struct {
	*Product
	Quantity int `db:"quantity"`
}

// OrderService ...
type OrderService interface {
	Get(id int) (*Order, error)
	GetByUser(id int) ([]*Order, error)
	Gets(first, limit int, order string) ([]*Order, error)

	Create(o *Order) error
	AddProducts(o *Order) error
	Update(o *Order, fields ...string) error
	Delete(id int) error
}
