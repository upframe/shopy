package mysql

import (
	"github.com/bruhs/shopy"
	"github.com/jmoiron/sqlx"
)

// UserService ...
type UserService struct {
	DB *sqlx.DB
}

var userMap = map[string]string{
	"ID":           "id",
	"FirstName":    "first_name",
	"LastName":     "last_name",
	"Email":        "email",
	"Address":      "address",
	"Invites":      "invites",
	"Credit":       "credit",
	"Confirmed":    "confirmed",
	"Admin":        "admin",
	"Referral":     "referral",
	"Referrer":     "referrer",
	"PasswordSalt": "password_salt",
	"PasswordHash": "password_hash",
	"Deactivated":  "deactivated",
}

// Get ...
func (s *UserService) Get(id int) (*shopy.User, error) {
	user := &shopy.User{}
	err := s.DB.Get(user, "SELECT * FROM users WHERE id=?", id)

	return user, err
}

// GetByEmail ...
func (s *UserService) GetByEmail(email string) (*shopy.User, error) {
	user := &shopy.User{}
	err := s.DB.Get(user, "SELECT * FROM users WHERE email=?", email)

	return user, err
}

// GetByReferral ...
func (s *UserService) GetByReferral(referral string) (*shopy.User, error) {
	user := &shopy.User{}
	err := s.DB.Get(user, "SELECT * FROM users WHERE referral=?", referral)

	return user, err
}

// Gets ...
func (s *UserService) Gets(first, limit int, order string) ([]*shopy.User, error) {
	users := []*shopy.User{}
	var err error

	order = fieldsToColumns(userMap, order)[0]

	if limit == 0 {
		err = s.DB.Select(&users, "SELECT * FROM users ORDER BY ?", order)
	} else {
		err = s.DB.Select(&users, "SELECT * FROM users ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)
	}

	return users, err
}

// Total ...
func (s *UserService) Total() (int, error) {
	return getTableCount(s.DB, "users")
}

// Create ...
func (s *UserService) Create(u *shopy.User) error {
	if u.ID != 0 {
		return nil
	}

	res, err := s.DB.NamedExec(insertQuery("users", getAllColumns(userMap)), u)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	u.ID = int(id)
	return err
}

// Update ...
func (s *UserService) Update(u *shopy.User, fields ...string) error {
	_, err := s.DB.NamedExec(updateQuery("users", "id", fieldsToColumns(userMap, fields...)), u)
	return err
}

// Delete ...
func (s *UserService) Delete(id int) error {
	u, err := s.Get(id)
	if err != nil {
		return err
	}

	u.Deactivated = true
	u.PasswordHash = ""
	u.PasswordSalt = ""
	return s.Update(u, "Deactivated", "PasswordHash", "PasswordSalt")
}
