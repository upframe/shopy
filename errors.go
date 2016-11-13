package fest

import "errors"

var (
	ErrAlreadyLoggedIn = errors.New("The user is already logged in.")
	ErrNotLoggedIn     = errors.New("The user is not logged in.")
)
