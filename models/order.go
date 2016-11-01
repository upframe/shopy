package models

// Order contains product order information
type Order struct {
	ID          int           `db:"id"`
	UserID      int           `db:"user_id"`
	PromocodeID NullInt64JSON `db:"promocode_id"`
	PayPalID    string        `db:"paypal_id"`
	Value       int           `db:"value"`
	Status      string        `db:"status"`
	Credits     int           `db:"credits"`
}

var orderColumns = []string{
	"id",
	"user_id",
	"promocode_id",
	"paypal_id",
	"value",
	"status",
	"credits",
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

// UserOrder will be used in /orders
type UserOrder struct {
	ID            int    `db:"id"`
	PromocodeID   int    `db:"promocode_id"`
	PromocodeCode string `db:"code"`
	Status        string `db:"status"`
	Value         int    `db:"value"`
	Products      []UserOrderProduct
}

// UserOrderProduct will be part of UserOrder
type UserOrderProduct struct {
	ID       int    `db:"product_id"`
	Name     string `db:"name"`
	Price    int    `db:"price"`
	Quantity int    `db:"quantity"`
}

// GetAllOrdersByUser gets all user orders
// TODO: format de sql to be easier to read; Do only two queries and process them
// with Go.
func GetAllOrdersByUser(user int) ([]UserOrder, error) {
	orders := []UserOrder{}
	err := db.Select(&orders, "SELECT o.id, pc.id as `promocode_id`, pc.code, o.status, o.value FROM upframe.orders AS o INNER JOIN upframe.promocodes AS pc ON o.promocode_id=pc.id AND o.user_id=?", user)

	for o := range orders {
		products := []UserOrderProduct{}
		db.Select(&products, "SELECT op.product_id, pd.name, pd.price, op.quantity FROM upframe.orders_products as op INNER JOIN upframe.products as pd ON op.product_id=pd.id AND op.order_id=?", orders[o].ID)
		orders[o].Products = products
	}

	return orders, err
}
