package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/upframe/fest"
)

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
	return &fest.Order{}, nil
}

// Gets ...
func (s *OrderService) Gets(first, limit int, order string) ([]*fest.Order, error) {
	return []*fest.Order{}, nil
}

// GetsWhere ...
func (s *OrderService) GetsWhere(first, limit int, order, where, sth string) ([]*fest.Order, error) {
	return []*fest.Order{}, nil
}

// Create ...
func (s *OrderService) Create(o *fest.Order) error {
	// TODO: add products and relationships
	return nil
}

// Update ...
func (s *OrderService) Update(o *fest.Order, fields ...string) error {
	// TODO: check if products or promocode to update
	return nil
}

// Delete ...
func (s *OrderService) Delete(id int) error {
	// TODO: just disable or change STATUS to 'CAnceled'
	return nil
}
