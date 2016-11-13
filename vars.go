package fest

import "errors"

const (
	passwordSaltBytes = 32
	passwordHashBytes = 64
)

var BaseAddress string

var (
	ErrAlreadyLoggedIn = errors.New("The user is already logged in.")
	ErrNotLoggedIn     = errors.New("The user is not logged in.")
)
