package mysql

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/upframe/fest"
)

var ordersMap = map[string]string{
	"ID":                    "o.id",
	"UserID":                "o.user_id",
	"PayPal":                "o.paypal",
	"Value":                 "o.value",
	"Status":                "o.Status",
	"Credits":               "o.Credits",
	"Promocode.ID":          "p.id",
	"Promocode.Code":        "p.code",
	"Promocode.Expires":     "p.expires",
	"Promocode.Discount":    "p.discount",
	"Promocode.Percentage":  "p.percentage",
	"Promocode.Deactivated": "p.deactivated",
}

type order struct {
	ID          int            `db:"id"`
	UserID      int            `db:"user_id"`
	PromocodeID fest.NullInt64 `db:"promocode_id"`
	PayPal      string         `db:"paypal"`
	Value       int            `db:"value"`
	Status      string         `db:"status"`
	Credits     int            `db:"credits"`
}

type orderProduct struct {
	ID        int64 `db:"id"`
	OrderID   int64 `db:"order_id"`
	ProductID int64 `db:"product_id"`
	Quantity  int   `db:"quantity"`
}

// OrderService ...
type OrderService struct {
	DB *sqlx.DB
}

// Get ...
func (s *OrderService) Get(id int) (*fest.Order, error) {
	orders, err := s.GetsWhere(0, 0, "ID", "ID", strconv.Itoa(id))
	if err != nil {
		return &fest.Order{}, err
	}

	if len(orders) == 0 {
		return &fest.Order{}, sql.ErrNoRows
	}

	return orders[0], nil
}

// Gets ...
func (s *OrderService) Gets(first, limit int, order string) ([]*fest.Order, error) {
	return s.GetsWhere(first, limit, order, "", "")
}

var orderBaseSelectQuery = "SELECT " +
	"o.id AS `order_id`," +
	"o.user_id AS `order_user`," +
	"o.paypal AS `order_paypal`," +
	"o.value AS `order_value`," +
	"o.status AS `order_status`," +
	"o.credits AS `order_credits`," +
	"pc.id AS `promocode_id`," +
	"pc.code AS `promocode_code`," +
	"pc.expires AS `promocode_expires`," +
	"pc.discount AS `promocode_discount`," +
	"pc.percentage AS `promocode_percentage`," +
	"pc.deactivated AS `promocode_deactivated` " +
	"FROM " +
	"orders AS o " +
	"LEFT JOIN " +
	"promocodes AS pc ON o.promocode_id = pc.id"

// GetsWhere ...
func (s *OrderService) GetsWhere(first, limit int, order, where, sth string) ([]*fest.Order, error) {
	var (
		rows *sql.Rows
		err  error
	)

	orders := []*fest.Order{}
	order = fieldsToColumns(ordersMap, order)[0]

	if where == "" {
		if limit == 0 {
			rows, err = s.DB.Query(orderBaseSelectQuery+" ORDER BY ?", order)
		} else {
			rows, err = s.DB.Query(orderBaseSelectQuery+" ORDER BY ? LIMIT ? OFFSET ?", limit, first)
		}
	} else {
		where = fieldsToColumns(ordersMap, where)[0]

		if limit == 0 {
			rows, err = s.DB.Query(orderBaseSelectQuery+" WHERE "+where+"=? ORDER BY ?", sth, order)
		} else {
			rows, err = s.DB.Query(orderBaseSelectQuery+" WHERE "+where+"=? ORDER BY ? LIMIT ? OFFSET ?", sth, limit, first)
		}
	}

	if err != nil {
		return orders, err
	}

	defer rows.Close()

	for rows.Next() {
		order := &fest.Order{Products: []*fest.OrderProduct{}}

		var (
			id          sql.NullInt64
			code        sql.NullString
			expires     sql.NullString
			discount    sql.NullInt64
			percentage  sql.NullBool
			deactivated sql.NullBool
		)

		err = rows.Scan(
			&order.ID, &order.UserID, &order.PayPal, &order.Value, &order.Status, &order.Credits,
			&id, &code, &expires, &discount, &percentage, &deactivated)
		if err != nil {
			return orders, err
		}

		if id.Valid {
			order.Promocode = &fest.Promocode{
				ID:          int(id.Int64),
				Code:        code.String,
				Discount:    int(discount.Int64),
				Percentage:  percentage.Bool,
				Deactivated: deactivated.Bool,
			}

			var t time.Time
			t, err = time.Parse(time.RFC3339, expires.String)
			if err != nil {
				return orders, err
			}

			order.Promocode.Expires = &t
		}

		var rowsps *sql.Rows
		rowsps, err = s.DB.Query("SELECT o.product_id, o.quantity, p.name FROM orders_products AS o INNER JOIN products AS p ON o.product_id = p.id WHERE o.order_id = ?", order.ID)
		defer rowsps.Close()

		for rowsps.Next() {
			prod := &fest.OrderProduct{}
			rowsps.Scan(&prod.ID, &prod.Quantity, &prod.Name)
			order.Products = append(order.Products, prod)
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return orders, err
	}

	return orders, nil
}

// Create ...
func (s *OrderService) Create(o *fest.Order) error {
	// TODO: add products and relationships
	return errors.New("Not implemented")
}

// Update ...
func (s *OrderService) Update(o *fest.Order, fields ...string) error {
	// TODO: check if products or promocode to update
	return errors.New("Not implemented")
}

// Delete ...
func (s *OrderService) Delete(id int) error {
	// TODO: just disable or change STATUS to 'CAnceled'
	return errors.New("Not implemented")
}
