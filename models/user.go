package models

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/scrypt"
)

const (
	passwordSaltBytes = 32
	passwordHashBytes = 64
)

// User contains the information of an User
type User struct {
	ID           int            `db:"id"`
	FirstName    string         `db:"first_name"`
	LastName     string         `db:"last_name"`
	Email        string         `db:"email"`
	Address      NullStringJSON `db:"address"`
	Invites      int            `db:"invites"`
	Credit       int            `db:"credit"`
	Confirmed    bool           `db:"confirmed"`
	Admin        bool           `db:"admin"`
	Referral     string         `db:"referral"`
	Referrer     NullInt64JSON  `db:"referrer"`
	PasswordSalt string         `db:"password_salt"`
	PasswordHash string         `db:"password_hash"`
	Deactivated  bool           `db:"deactivated"`
}

var userColumns = []string{
	"id",
	"first_name",
	"last_name",
	"email",
	"address",
	"invites",
	"credit",
	"confirmed",
	"admin",
	"referral",
	"referrer",
	"password_salt",
	"password_hash",
	"deactivated",
}

// Insert inserts the current User struct into the database and returns an error
// if something goes wrong.
func (u User) Insert() (int64, error) {
	if u.ID != 0 {
		return 0, nil
	}

	res, err := db.NamedExec(insertQuery("users", userColumns), u)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Update updates the current User struct into the database
func (u User) Update(fields ...string) error {
	if fields[0] == UpdateAll {
		fields = userColumns
	}

	_, err := db.NamedExec(updateQuery("users", "id", fields), u)
	return err
}

// Deactivate deletes a user from the database using his id
func (u *User) Deactivate() error {
	u.Deactivated = true
	return u.Update("deactivated")
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

// CheckPassword checks if the password of the user is correct
func (u User) CheckPassword(password string) (bool, error) {
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

// GenerateReferralHash generates a new referrer hash for a new user
func (u *User) GenerateReferralHash() {
	if u.Referral != "" {
		return
	}

	u.Referral = UniqueHash(u.Email)
}

// GetUserByID retrieves a user from the database using its ID
func GetUserByID(id int) (Generic, error) {
	user := User{}
	err := db.Get(&user, "SELECT * FROM users WHERE id=?", id)
	return &user, err
}

// GetUserByEmail retrieves a user from the database using its email
func GetUserByEmail(email string) (*User, error) {
	user := User{}
	err := db.Get(&user, "SELECT * FROM users WHERE email=?", email)
	return &user, err
}

// GetUserByReferral retrieves a user from the database using its email
func GetUserByReferral(referral string) (*User, error) {
	user := User{}
	err := db.Get(&user, "SELECT * FROM users WHERE referral=?", referral)
	return &user, err
}

// GetUsers retrieves more than one user from the databse
func GetUsers(first, limit int, order string) ([]Generic, error) {
	users := []User{}
	err := db.Select(&users, "SELECT * FROM users ORDER BY ? LIMIT ? OFFSET ?", order, limit, first)

	generics := make([]Generic, len(users))
	for i := range users {
		generics[i] = &users[i]
	}
	return generics, err
}
