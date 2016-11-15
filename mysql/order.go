package mysql

import "github.com/jmoiron/sqlx"

// OrderService ...
type OrderService struct {
	DB *sqlx.DB
}
