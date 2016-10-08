package models

import "github.com/gorilla/sessions"

type Session struct {
	*sessions.Session
	User *User
}

func (s Session) IsLoggedIn() bool {
	switch s.Values["IsLoggedIn"].(type) {
	case bool:
		return s.Values["IsLoggedIn"].(bool)
	}

	return false
}

// IsAdmin checks if an user is admin
func (s Session) IsAdmin() bool {
	if !s.IsLoggedIn() {
		return false
	}

	return s.User.Admin
}
