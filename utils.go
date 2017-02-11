package shopy

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

// UniqueHash returns a SHA256 hash based on the string and on the
// current time
func UniqueHash(phrase string) string {
	data := phrase + time.Now().Format(time.ANSIC)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// DisplayCents ...
func DisplayCents(cents int) string {
	price := strconv.Itoa(cents)

	if len(price) == 1 {
		price = "0.0" + price
	} else if len(price) == 2 {
		price = "0." + price
	} else {
		cents := price[len(price)-2:]
		price = price[0:len(price)-2] + "." + cents
	}

	return price
}
