package mysql

import (
	"github.com/upframe/fest"
)

// ProductService ...
type ProductService struct{}

var productMap = map[string]string{
	"ID":          "id",
	"Name":        "name",
	"Description": "description",
	"Price":       "price",
	"Picture":     "picture",
	"Deactivated": "deactivated",
}

// Get ...
func (s *ProductService) Get(id int) (*fest.Product, error) {
	product := &fest.Product{}
	err := db.Get(product, "SELECT * FROM products WHERE id=?", id)

	return product, err
}

// Gets ...
func (s *ProductService) Gets(first, limit int, order string) ([]*fest.Product, error) {
	products := []*fest.Product{}
	var err error

	if limit == 0 {
		err = db.Select(&products, "SELECT * FROM products ORDER BY ?", order)
	} else {
		err = db.Select(&products, "SELECT * FROM products ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)
	}

	return products, err
}

// GetsWhere ...
func (s *ProductService) GetsWhere(first, limit int, order, where, sth string) ([]*fest.Product, error) {
	products := []*fest.Product{}
	var err error

	if limit == 0 {
		err = db.Select(&products, "SELECT * FROM products WHERE "+where+"=? ORDER BY ?", sth, order)
	} else {
		err = db.Select(&products, "SELECT * FROM products WHERE "+where+"=? ORDER BY ? LIMIT ? OFFSET ?", sth, order, limit, first)
	}

	return products, err
}

// GetsWhereIn ...
func (s *ProductService) GetsWhereIn(first, limit int, order, where, in string) ([]*fest.Product, error) {
	products := []*fest.Product{}
	var err error

	if limit == 0 {
		err = db.Select(&products, "SELECT * FROM products WHERE "+where+" IN ? ORDER BY ?", in, order)
	} else {
		err = db.Select(&products, "SELECT * FROM products WHERE "+where+" IN ? ORDER BY ? LIMIT ? OFFSET ?", in, order, limit, first)
	}

	return products, err
}

// Create ...
func (s *ProductService) Create(p *fest.Product) error {
	if p.ID != 0 {
		return nil
	}

	res, err := db.NamedExec(insertQuery("products", getAllColumns(productMap)), p)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	p.ID = int(id)
	return err
}

// Update ...
func (s *ProductService) Update(p *fest.Product, fields ...string) error {
	_, err := db.NamedExec(updateQuery("products", "id", fieldsToColumns(productMap, fields...)), p)
	return err
}

// Delete ...
func (s *ProductService) Delete(p *fest.Product) error {
	p.Deactivated = true
	return s.Update(p, "Deactivated")
}
