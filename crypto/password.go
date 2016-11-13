package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/upframe/fest"
	"golang.org/x/crypto/scrypt"
)

const (
	passwordSaltBytes = 32
	passwordHashBytes = 64
)

// CheckPassword checks if the password of the user is correct
func CheckPassword(u *fest.User, password string) (bool, error) {
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
func SetPassword(u *fest.User, password string) error {
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
