package random

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// UniqueHash returns a SHA256 hash based on the string and on the
// current time
func UniqueHash(phrase string) string {
	data := phrase + time.Now().Format(time.ANSIC)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
