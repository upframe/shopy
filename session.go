package fest

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Session ...
type Session struct {
	*sessions.Session
	User *User
}

// IsLoggedIn checks if the user is logged in
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

// SessionService ...
type SessionService interface {
	Session(w http.ResponseWriter, r *http.Request) (*Session, error)
	Cart(s *Session) (*Cart, error)
}
