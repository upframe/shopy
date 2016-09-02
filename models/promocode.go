package models

import "time"

// Promocode contains promocodes informations
type Promocode struct {
	ID          int        `db:"id"`
	Code        string     `db:"code"`
	Expires     *time.Time `db:"expires"`
	Discount    int        `db:"discount"`
	Deactivated bool       `db:"deactivated"`
}

// Insert inserts an order into the database
func (p Promocode) Insert() error {
	if p.ID != 0 {
		return nil
	}

	_, err := db.NamedExec(
		`INSERT INTO promocodes
								(id,
									code,
									expires,
									discount,
									deactivated
			VALUES 		(:id,
									:code,
									:expires,
									:discount,
									:deactivated)`, p)

	return err
}

// Update updates an order from the database
func (p Promocode) Update(fields ...string) error {
	_, err := db.NamedExec(updateQuery("promocodes", "id", fields), p)
	return err
}

// Deactivate deactivates a promocode
func (p *Promocode) Deactivate() error {
	p.Deactivated = true
	return p.Update("deactivated")
}

// GetPromocode pulls out an order from the database
func GetPromocode(id int) (Generic, error) {
	promocode := &Promocode{}
	err := db.Get(promocode, "SELECT * FROM promocodes WHERE id=?", id)

	return promocode, err
}

// GetPromocodes does something that I don't actually know
func GetPromocodes(first, limit int) ([]Generic, error) {
	promocodes := []Promocode{}
	err := db.Select(&promocodes, "SELECT * FROM promocodes ORDER BY id LIMIT ? OFFSET ?", limit, first)
	//fmt.Println(promocodes)
	generics := make([]Generic, len(promocodes))
	for i := range promocodes {
		generics[i] = &promocodes[i]
	}
	return generics, err
}
