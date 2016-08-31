package models

// Product contains products informations
type Product struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Price       int    `db:"price"`
	Picture     string `db:"picture"`
	Deactivated bool   `db:"deactivated"`
}
