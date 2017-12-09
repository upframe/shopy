package mysql

import (
	"fmt"

	"github.com/upframe/shopy"
	"github.com/jmoiron/sqlx"
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
	"Usage":       "usage",
}

// Get ...
func (s *PromocodeService) Get(id int) (*shopy.Promocode, error) {
	promocode := &shopy.Promocode{}
	err := s.DB.Get(promocode, "SELECT * FROM promocodes WHERE id=?", id)

	return promocode, err
}

// GetByCode ...
func (s *PromocodeService) GetByCode(code string) (*shopy.Promocode, error) {
	promocode := &shopy.Promocode{}
	err := s.DB.Get(promocode, "SELECT * FROM promocodes WHERE code=?", code)

	return promocode, err
}

// Gets ...
func (s *PromocodeService) Gets(first, limit int, order string) ([]*shopy.Promocode, error) {
	promocodes := []*shopy.Promocode{}
	var err error

	order = fieldsToColumns(promocodeMap, order)[0]

	if limit == 0 {
		err = s.DB.Select(&promocodes, "SELECT * FROM promocodes ORDER BY ?", order)
	} else {
		err = s.DB.Select(&promocodes, "SELECT * FROM promocodes ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)
	}

	return promocodes, err
}

// Total ...
func (s *PromocodeService) Total() (int, error) {
	return getTableCount(s.DB, "promocodes")
}

// Create ...
func (s *PromocodeService) Create(p *shopy.Promocode) error {
	if p.ID != 0 {
		return nil
	}

	fmt.Println(insertQuery("promocodes", getAllColumns(promocodeMap)))

	res, err := s.DB.NamedExec(insertQuery("promocodes", getAllColumns(promocodeMap)), p)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	p.ID = int(id)
	return err
}

// Update ...
func (s *PromocodeService) Update(p *shopy.Promocode, fields ...string) error {
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
