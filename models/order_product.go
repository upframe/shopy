package models

// OrderProduct contains the relationship between a product and an order
type OrderProduct struct {
  ID           int           `db:"id"` // TODO: remove id or not?
  OrderID      int           `db:"order_id"`
  ProductID    int           `db:"product_id"`
}

var orderProductColumns = []string{
	"id",
	"order_id",
	"product_id",
}

// Insert inserts an order-product into the database
func (op OrderProduct) Insert() (int64, error) {
	if op.ID != 0 {
		return 0, nil
	}

	res, err := db.NamedExec(insertQuery("orders_products", orderProductColumns), op)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Update updates an order from the database
func (op OrderProduct) Update(fields ...string) error {
	if fields[0] == UpdateAll {
		fields = orderProductColumns
	}

	_, err := db.NamedExec(updateQuery("orders_products", "id", fields), op)
	return err
}

// Deactivate deactivates an order
func (op OrderProduct) Deactivate() error {
	return nil
}

// NOTE: incorporate this on order? Get Order returns  JOIN? Or Order should have a method like GetProducts that returns an array of OrderProduct?

