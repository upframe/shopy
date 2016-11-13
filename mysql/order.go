package mysql

import (
	"github.com/upframe/fest"
)

// OrderService ...
type OrderService struct{}

var orderColumns = []string{
	"id",
	"user_id",
	"promocode_id",
	"paypal_id",
	"value",
	"status",
	"credits",
}

// Order ...
func (s *OrderService) Order(id int) (*fest.Order, error) {
	order := &fest.Order{}
	err := db.Get(order, "SELECT * FROM orders WHERE id=?", id)

	return order, err
}

// Orders ...
func (s *OrderService) Orders(first, limit int, order string) ([]*fest.Order, error) {
	orders := []*fest.Order{}
	var err error

	if limit == 0 {
		err = db.Select(&orders, "SELECT * FROM orders ORDER BY ?", order)
	} else {
		err = db.Select(&orders, "SELECT * FROM orders ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)
	}

	return orders, err
}

// UserOrders ...
func (s *OrderService) UserOrders(u *fest.User) ([]*fest.Order, error) {
	/*orders := []UserOrder{}
	err := db.Select(&orders, "SELECT o.id, pc.id as `promocode_id`, pc.code, o.status, o.value FROM upframe.orders AS o INNER JOIN upframe.promocodes AS pc ON o.promocode_id=pc.id AND o.user_id=?", user)

	for o := range orders {
		products := []UserOrderProduct{}
		db.Select(&products, "SELECT op.product_id, pd.name, pd.price, op.quantity FROM upframe.orders_products as op INNER JOIN upframe.products as pd ON op.product_id=pd.id AND op.order_id=?", orders[o].ID)
		orders[o].Products = products
	}

	return orders, err*/

	return []*fest.Order{}, nil
}

// CreateOrder ...
func (s *OrderService) CreateOrder(o *fest.Order) error {
	if o.ID != 0 {
		return nil
	}

	res, err := db.NamedExec(insertQuery("orders", orderColumns), o)
	if err != nil {
		return err
	}

	o.ID, err = res.LastInsertId()
	return err
}

// UpdateOrder ...
func (s *OrderService) UpdateOrder(o *fest.Order, fields ...string) error {
	if fields[0] == UpdateAll {
		fields = orderColumns
	}

	_, err := db.NamedExec(updateQuery("orders", "id", fields), o)
	return err
}

// DeleteOrder ...
func (s *OrderService) DeleteOrder(o *fest.Order) error {
	return nil
}
