package models

import "database/sql"

// Order contains product order information
type Order struct {
	ID          int           `db:"id"`
	UserID      int           `db:"user_id"`
	ProductID   int           `db:"product_id"`
	PromocodeID sql.NullInt64 `db:"promocode_id"`
	Value       int           `db:"value"`
}
