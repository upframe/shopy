package upframe

import "time"

// User is an user
type User struct {
	ID              int
	FullName        string
	EmailAddress    string
	PhysicalAddress string
	IsConfirmed     bool
	CreatedAt       *time.Time
	PasswordHash    string
	PasswordSalt    string
}
