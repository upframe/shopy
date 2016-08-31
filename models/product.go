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

// Insert inserts an order into the database
func (p Product) Insert() error {
	if p.ID != 0 {
		return nil
	}

	_, err := db.NamedExec(
		`INSERT INTO products
								(id,
									name,
									description,
									price,
									picture,
									deactivated)
			VALUES 		(:id,
									:name,
									:description,
									:price,
									:picture,
									:deactivated)`, p)

	return err
}

// Update updates an order from the database
func (p Product) Update(fields ...string) error {
	_, err := db.NamedExec(updateQuery("products", "id", fields), p)
	return err
}

// GetProducts pulls out an order from the database
func GetProduct(id int) (*Product, error) {
	product := &Product{}
	err := db.Get(product, "SELECT * FROM products WHERE id=?", id)

	return product, err
}

// GetProducts does something that I don't actually know
func GetProducts(first, limit int) ([]*Product, error) {
	products := []*Product{}
	err := db.Select(products, "SELECT * FROM products ORDER BY id LIMIT ? OFFSET ?", limit, first)

	return products, err
}
