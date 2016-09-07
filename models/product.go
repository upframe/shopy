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

var productColumns = []string{
	"id",
	"name",
	"description",
	"price",
	"picture",
	"deactivated",
}

// Insert inserts a product into the database
func (p Product) Insert() error {
	if p.ID != 0 {
		return nil
	}

	_, err := db.NamedExec(insertQuery("products", productColumns), p)
	return err
}

// Update updates a product from the database
func (p Product) Update(fields ...string) error {
	if fields[0] == UpdateAll {
		fields = productColumns
	}

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
	var err error

	if limit == 0 {
		err = db.Select(&products, "SELECT * FROM products ORDER BY ?", order)
	} else {
		err = db.Select(&products, "SELECT * FROM products ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)
	}

	generics := make([]Generic, len(products))
	for i := range products {
		generics[i] = &products[i]
	}

	return generics, err
}

// GetProductsWhere retrives products from the database
func GetProductsWhere(first, limit int, order string, where string, sth string) ([]Generic, error) {
	products := []Product{}
	var err error

	if limit == 0 {
		err = db.Select(&products, "SELECT * FROM products WHERE "+where+"=? ORDER BY ?", sth, order)
	} else {
		err = db.Select(&products, "SELECT * FROM products WHERE "+where+"=? ORDER BY ? LIMIT ? OFFSET ?", sth, order, limit, first)
	}

	generics := make([]Generic, len(products))
	for i := range products {
		generics[i] = &products[i]
	}

	return generics, err
}
