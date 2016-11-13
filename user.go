package fest

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/scrypt"
)

// User contains the information of an User
type User struct {
	ID           int        `db:"id"`
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

// CheckPassword checks if the password of the user is correct
func (u *User) CheckPassword(password string) (bool, error) {
	// Decodes the hexadecimal salt into a []byte
	salt, err := hex.DecodeString(u.PasswordSalt)
	if err != nil {
		return false, err
	}

	// Makes an hash from the password and salt
	hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, passwordHashBytes)
	if err != nil {
		return false, err
	}

	return (hex.EncodeToString(hash) == u.PasswordHash), nil
}

// SetPassword generates the salt and the hash of the user password
func (u *User) SetPassword(password string) error {
	// Generates a random salt
	salt := make([]byte, passwordSaltBytes)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return err
	}

	// Creates the hash from the salt and password
	hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, passwordHashBytes)
	if err != nil {
		return err
	}

	// Stores the hash and salt in hexadecimal form
	u.PasswordSalt = hex.EncodeToString(salt)
	u.PasswordHash = hex.EncodeToString(hash)
	return nil
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
