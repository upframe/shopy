package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/upframe/fest"
)

// PromocodeService ...
type PromocodeService struct {
	DB *sqlx.DB
}

var promocodeMap = map[string]string{
	"ID":          "id",
	"Code":        "code",
	"Expires":     "expires",
	"Discount":    "discount",
	"Percentage":  "percentage",
	"Deactivated": "deactivated",
}

// Get ...
func (s *PromocodeService) Get(id int) (*fest.Promocode, error) {
	promocode := &fest.Promocode{}
	err := s.DB.Get(promocode, "SELECT * FROM promocodes WHERE id=?", id)

	return promocode, err
}

// GetByCode ...
func (s *PromocodeService) GetByCode(code string) (*fest.Promocode, error) {
	promocode := &fest.Promocode{}
	err := s.DB.Get(promocode, "SELECT * FROM promocodes WHERE code=?", code)

	return promocode, err
}

// Gets ...
func (s *PromocodeService) Gets(first, limit int, order string) ([]*fest.Promocode, error) {
	promocodes := []*fest.Promocode{}
	var err error

	if limit == 0 {
		err = s.DB.Select(&promocodes, "SELECT * FROM promocodes ORDER BY ?", order)
	} else {
		err = s.DB.Select(&promocodes, "SELECT * FROM promocodes ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)
	}

	return promocodes, err
}

// Create ...
func (s *PromocodeService) Create(p *fest.Promocode) error {
	if p.ID != 0 {
		return nil
	}

	res, err := s.DB.NamedExec(insertQuery("promocodes", getAllColumns(promocodeMap)), p)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	p.ID = int(id)
	return err
}

// Update ...
func (s *PromocodeService) Update(p *fest.Promocode, fields ...string) error {
	_, err := s.DB.NamedExec(updateQuery("promocodes", "id", fieldsToColumns(promocodeMap, fields...)), p)
	return err
}

// Delete ...
func (s *PromocodeService) Delete(id int) error {
	p, err := s.Get(id)
	if err != nil {
		return err
	}

	p.Deactivated = true
	return s.Update(p, "Deactivated")
}
