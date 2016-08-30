package models

import "time"

// Promocode contains promocodes informations
type Promocode struct {
	ID       int        `db:"id"`
	Code     string     `db:"code"`
	Expires  *time.Time `db:"expires"`
	Discount int        `db:"discount"`
}
