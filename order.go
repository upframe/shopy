package shopy

const (
	// OrderCanceled represents the order canceled status.
	OrderCanceled = -1
	// OrderPaymentWaiting represents an waiting payment status.
	OrderPaymentWaiting = 0
	// OrderPaymentDone represents the payment done status.
	OrderPaymentDone = 1
	// OrderPaymentFailed represents the failed payment status.
	OrderPaymentFailed = 2
	// OrderUnfulfilled represents a unfulfilled order.
	OrderUnfulfilled = 0
	// OrderFulfilled represents a fulfilled order.
	OrderFulfilled = 1
)

// Order ...
type Order struct {
	ID                int
	PayPal            string
	Value             int
	PaymentStatus     int16
	FulfillmentStatus int16
	Credits           int
	User              *User
	Promocode         *Promocode
	Products          []*OrderProduct
}

// PaymentStatusText returns the text corresponding to the status variable.
func (o *Order) PaymentStatusText() string {
	switch o.PaymentStatus {
	case OrderPaymentWaiting:
		return "Waiting"
	case OrderPaymentDone:
		return "Done"
	case OrderPaymentFailed:
		return "Failed"
	case OrderCanceled:
		return "Canceled"
	}

	return "Unknown"
}

// FulfillmentStatusText returns the text corresponding to the status variable.
func (o *Order) FulfillmentStatusText() string {
	switch o.FulfillmentStatus {
	case OrderFulfilled:
		return "Fulfilled"
	case OrderUnfulfilled:
		return "Unfulfilled"
	case OrderCanceled:
		return "Canceled"
	}

	return "Unknown"
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

	Total() (int, error)
	Create(o *Order) error
	Update(o *Order, fields ...string) error
	Delete(id int) error
}
