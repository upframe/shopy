package fest

import "time"

// Promocode ...
type Promocode struct {
	ID          int        `db:"id"`
	Code        string     `db:"code"`
	Expires     *time.Time `db:"expires"`
	Discount    int        `db:"discount"`
	Percentage  bool       `db:"percentage"`
	Deactivated bool       `db:"deactivated"`
}

// PromocodeService ...
type PromocodeService interface {
	Get(id int) (*Promocode, error)
	GetByCode(code string) (*Promocode, error)
	Gets(first, limit int, order string) ([]*Promocode, error)

	Total() (int, error)
	Create(p *Promocode) error
	Update(p *Promocode, fields ...string) error
	Delete(id int) error
}
