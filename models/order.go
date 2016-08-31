package models

import "database/sql"

// Order contains product order information
type Order struct {
	ID          int           `db:"id"`
	UserID      int           `db:"user_id"`
	ProductID   int           `db:"product_id"`
	PromocodeID sql.NullInt64 `db:"promocode_id"`
	Value       int           `db:"value"`
	Status      int           `db:"status"`
}

// Insert inserts an order into the database
func (o Order) Insert() error {
	if o.ID != 0 {
		return nil
	}

	_, err := db.NamedExec(
		`INSERT INTO orders
								(id,
									user_id,
									product_id,
									promocode_id
									value)
			VALUES 		(:id,
									:user_id,
									:product_id,
									:promocode_id
									:value)`, o)

	return err
}

// Update updates an order from the database
func (o Order) Update(fields ...string) error {
	_, err := db.NamedExec(updateQuery("orders", "id", fields), o)
	return err
}

// GetOrder pulls out an order from the database
func GetOrder(id int) (*Order, error) {
	order := &Order{}
	err := db.Get(order, "SELECT * FROM orders WHERE id=?", id)

	return order, err
}

// GetOrders does something that I don't actually know
func GetOrders(first, limit int) (*[]Order, error) {
	orders := &[]Order{}
	err := db.Select(orders, "SELECT * FROM orders ORDER BY id LIMIT ? OFFSET ?", limit, first)

	return orders, err
}
