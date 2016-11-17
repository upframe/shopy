package fest

const (
	OrderWaitingPayment = "pending"
	OrderApproved       = "approved"
	OrderCreated        = "created"
	OrderFailed         = "failed"
	OrderCanceled       = "canceled"
)

// Order ...
type Order struct {
	ID        int
	PayPal    string
	Value     int
	Status    string
	Credits   int
	User      *User
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
	GetByPayPal(token string) (*Order, error)
	Gets(first, limit int, order string) ([]*Order, error)
	GetsWhere(first, limit int, order, where, sth string) ([]*Order, error)

	Create(o *Order) error
	Update(o *Order, fields ...string) error
	Delete(id int) error
}
