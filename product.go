package fest

// Product ...
type Product struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Price       int    `db:"price"`
	Picture     string `db:"picture"`
	Deactivated bool   `db:"deactivated"`
}

// ProductService ...
type ProductService interface {
	Get(id int) (*Product, error)
	Gets(first, limit int, order string) ([]*Product, error)
	GetsWhere(first, limit int, order, where, sth string) ([]*Product, error)
	GetsWhereIn(first, limit int, order, where, in string) ([]*Product, error)

	Total() (int, error)
	Create(p *Product) error
	Update(p *Product, fields ...string) error
	Delete(id int) error
}
