package models

import "database/sql"

// User contains the information of an User
type User struct {
	FirstName       string
	LastName        string
	EmailAddress    string
	PhysicalAddress sql.NullString
}

// Update updates the current User struct into the database
func (u User) Update() error {
	return nil
}

// Insert inserts the current User struct into the database
func (u User) Insert() error {
	return nil
}

// DeleteUser deletes a user from the database using its email
func DeleteUser(email string) error {
	return nil
}

// GetUser retrieves a user from the database using its email
func GetUser(email string) (*User, error) {
	return &User{}, nil
}
