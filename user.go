package fest

// User contains the information of an User
type User struct {
	ID           int      `db:"id"`
	FirstName    string     `db:"first_name"`
	LastName     string     `db:"last_name"`
	Email        string     `db:"email"`
	Address      NullString `db:"address"`
	Invites      int        `db:"invites"`
	Credit       int        `db:"credit"`
	Confirmed    bool       `db:"confirmed"`
	Admin        bool       `db:"admin"`
	Referral     string     `db:"referral"`
	Referrer     NullInt64  `db:"referrer"`
	PasswordSalt string     `db:"password_salt"`
	PasswordHash string     `db:"password_hash"`
	Deactivated  bool       `db:"deactivated"`
}

// UserService ...
type UserService interface {
	Get(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByReferral(referral string) (*User, error)
	Gets(first, limit int, order string) ([]*User, error)

	Create(u *User) (int, error)
	Update(u *User, fields ...string) error
	Delete(id int) error
}
