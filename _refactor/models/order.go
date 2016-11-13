package models

// GetAllOrdersByUser gets all user orders
// TODO: format de sql to be easier to read; Do only two queries and process them;
// BUG: this do not retrieve promocodeless orders...
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
