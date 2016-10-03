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

var promocodeColumns = []string{
	"id",
	"code",
	"expires",
	"discount",
	"deactivated",
}

// Insert inserts an order into the database
func (p Promocode) Insert() (int64, error) {
	if p.ID != 0 {
		return 0, nil
	}

	res, err := db.NamedExec(insertQuery("promocodes", promocodeColumns), p)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Update updates an order from the database
func (p Promocode) Update(fields ...string) error {
	if fields[0] == UpdateAll {
		fields = promocodeColumns
	}

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

// GetPromocodeByCode gets a promocode from the database using the code
func GetPromocodeByCode(code string) (Generic, error) {
	promocode := &Promocode{}
	err := db.Get(promocode, "SELECT * FROM promocodes WHERE code=?", code)

	return promocode, err
}

// GetPromocodes does something that I don't actually know
func GetPromocodes(first, limit int, order string) ([]Generic, error) {
	promocodes := []Promocode{}
	err := db.Select(&promocodes, "SELECT * FROM promocodes ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)
	//fmt.Println(promocodes)
	generics := make([]Generic, len(promocodes))
	for i := range promocodes {
		generics[i] = &promocodes[i]
	}
	return generics, err
}
