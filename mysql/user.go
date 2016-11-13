package mysql

import (
	"github.com/upframe/fest"
)

// UserService ...
type UserService struct{}

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
func (s *UserService) Get(id int) (*fest.User, error) {
	user := &fest.User{}
	err := db.Get(user, "SELECT * FROM users WHERE id=?", id)

	return user, err
}

// GetByEmail ...
func (s *UserService) GetByEmail(email string) (*fest.User, error) {
	user := &fest.User{}
	err := db.Get(user, "SELECT * FROM users WHERE email=?", email)

	return user, err
}

// GetByReferral ...
func (s *UserService) GetByReferral(referral string) (*fest.User, error) {
	user := &fest.User{}
	err := db.Get(user, "SELECT * FROM users WHERE referral=?", referral)

	return user, err
}

// Gets ...
func (s *UserService) Gets(first, limit int, order string) ([]*fest.User, error) {
	users := []*fest.User{}
	var err error

	if limit == 0 {
		err = db.Select(&users, "SELECT * FROM users ORDER BY ?", order)
	} else {
		err = db.Select(&users, "SELECT * FROM users ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)
	}

	return users, err
}

// Create ...
func (s *UserService) Create(u *fest.User) error {
	if u.ID != 0 {
		return nil
	}

	res, err := db.NamedExec(insertQuery("users", getAllColumns(userMap)), u)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	u.ID = int(id)
	return err
}

// Update ...
func (s *UserService) Update(u *fest.User, fields ...string) error {
	_, err := db.NamedExec(updateQuery("users", "id", fieldsToColumns(userMap, fields...)), u)
	return err
}

// Delete ...
func (s *UserService) Delete(id int) error {
	u, err := s.Get(id)
	if err != nil {
		return err
	}

	u.Deactivated = true
	return s.Update(u, "Deactivated")
}
