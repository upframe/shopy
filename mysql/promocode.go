package mysql

import (
	"github.com/upframe/fest"
)

// PromocodeService ...
type PromocodeService struct{}

var promocodeColumns = []string{
	"id",
	"code",
	"expires",
	"discount",
	"percentage",
	"deactivated",
}

// Promocode ...
func (s *PromocodeService) Promocode(id int) (*fest.Promocode, error) {
	promocode := &fest.Promocode{}
	err := db.Get(promocode, "SELECT * FROM promocodes WHERE id=?", id)

	return promocode, err
}

// PromocodeByCode ...
func (s *PromocodeService) PromocodeByCode(code string) (*fest.Promocode, error) {
	promocode := &fest.Promocode{}
	err := db.Get(promocode, "SELECT * FROM promocodes WHERE code=?", code)

	return promocode, err
}

// Promocodes ...
func (s *PromocodeService) Promocodes(first, limit int, order string) ([]*fest.Promocode, error) {
	promocodes := []*fest.Promocode{}
	err := db.Select(&promocodes, "SELECT * FROM promocodes ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)

	return promocodes, err
}

// CreatePromocode ...
func (s *PromocodeService) CreatePromocode(p *fest.Promocode) error {
	if p.ID != 0 {
		return nil
	}

	res, err := db.NamedExec(insertQuery("promocodes", promocodeColumns), p)
	if err != nil {
		return err
	}

	p.ID, err = res.LastInsertId()
	return err
}

// UpdatePromocode ...
func (s *PromocodeService) UpdatePromocode(p *fest.Promocode, fields ...string) error {
	if fields[0] == UpdateAll {
		fields = promocodeColumns
	}

	_, err := db.NamedExec(updateQuery("promocodes", "id", fields), p)
	return err
}

// DeletePromocode ...
func (s *PromocodeService) DeletePromocode(p *fest.Promocode) error {
	p.Deactivated = true
	return s.UpdatePromocode(p, "deactivated")
}
