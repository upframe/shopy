package upframe

import "time"

type Order struct {
	ID        int
	User      *User
	Product   *Product
	Promocode *Promocode
	Value     int
	CreatedAt *time.Time
}

type Product struct {
	ID          int
	Name        string
	Description string
	Price       int
	Picture     string
	CreatedAt   *time.Tine
}

type Promocode struct {
	ID        int
	Code      string
	CreatedAt *time.Time
	Validity  *time.Time
	Discount  int
}
