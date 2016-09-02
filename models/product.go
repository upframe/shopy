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

// Insert inserts a product into the database
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

// Update updates a product from the database
func (p Product) Update(fields ...string) error {
	_, err := db.NamedExec(updateQuery("products", "id", fields), p)
	return err
}

// Deactivate deactivates a product (changes its visibility)
func (p Product) Deactivate() error {
	p.Deactivated = true
	return p.Update("deactivated")
}

// GetProduct retrieves a product from the database
func GetProduct(id int) (Generic, error) {
	product := &Product{}
	err := db.Get(product, "SELECT * FROM products WHERE id=?", id)

	return product, err
}

// GetProducts retrives products from the database
func GetProducts(first, limit int, order string) ([]Generic, error) {
	products := []Product{}
	err := db.Select(&products, "SELECT * FROM products ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)

	generics := make([]Generic, len(products))
	for i := range products {
		generics[i] = &products[i]
	}

	return generics, err
}
