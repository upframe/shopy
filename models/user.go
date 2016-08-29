package models

import "database/sql"

// User contains the information of an User
type User struct {
	FirstName    string         `db:"first_name"`
	LastName     string         `db:"last_name"`
	Email        string         `db:"email"`
	Address      sql.NullString `db:"address"`
	Invites      int            `db:"invites"`
	Credit       int            `db:"credit"`
	Confirmed    bool           `db:"confirmed"`
	ReferrerHash string         `db:"referrer_hash"`
	ReferredBy   int            `db:"referred_by"`
	PasswordSalt string         `db:"password_salt"`
	PasswordHash string         `db:"password_hash"`
}

// Update updates the current User struct into the database
func (u User) Update() error {
	return nil
}

// Insert inserts the current User struct into the database
func (u User) Insert() error {
	return nil
}

// CheckPassword checks if the password of the user is correct
func (u User) CheckPassword(password string) bool {
	return false
}

// DeleteUser deletes a user from the database using its email
func DeleteUser(email string) error {
	return nil
}

// GetUser retrieves a user from the database using its email
func GetUser(email string) (*User, error) {
	return &User{}, nil
}
