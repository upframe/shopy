package models

// Order contains product order information
type Order struct {
	ID          int           `db:"id"`
	UserID      int           `db:"user_id"`
	PromocodeID NullInt64JSON `db:"promocode_id"`
	Value       int           `db:"value"`
	Status      int           `db:"status"`
	PayPalID    string        `db:"paypal_id"`
}

var orderColumns = []string{
	"id",
	"user_id",
	"promocode_id",
	"value",
	"paypal_id",
}

// Insert inserts an order into the database
func (o Order) Insert() (int64, error) {
	if o.ID != 0 {
		return 0, nil
	}

	res, err := db.NamedExec(insertQuery("orders", orderColumns), o)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Update updates an order from the database
func (o Order) Update(fields ...string) error {
	if fields[0] == UpdateAll {
		fields = orderColumns
	}

	_, err := db.NamedExec(updateQuery("orders", "id", fields), o)
	return err
}

// Deactivate deactivates an order
func (o Order) Deactivate() error {
	return nil
}

// GetOrder pulls out an order from the database
func GetOrder(id int) (Generic, error) {
	order := &Order{}
	err := db.Get(order, "SELECT * FROM orders WHERE id=?", id)

	return order, err
}

// GetOrders does something that I don't actually know
func GetOrders(first, limit int, order string) ([]Generic, error) {
	orders := []Order{}
	err := db.Select(&orders, "SELECT * FROM orders ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)

	generics := make([]Generic, len(orders))
	for i := range orders {
		generics[i] = &orders[i]
	}

	return generics, err
}

// GetAllOrdersByUser gets all user orders
func GetAllOrdersByUser(user int) ([]Generic, error) {
	orders := []Order{}
	err := db.Select(&orders, "SELECT * FROM orders WHERE user_id =?", user)

	generics := make([]Generic, len(orders))
	for i := range orders {
		generics[i] = orders[i]
	}

	return generics, err
}
