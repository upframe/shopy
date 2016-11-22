package fest

import "net/http"

// Session ...
type Session struct {
	Logged bool
	User   *User
}

// SessionService ...
type SessionService interface {
	Save(w http.ResponseWriter, sess *Session) error
	Get(w http.ResponseWriter, r *http.Request) (*Session, error)
	Reset(w http.ResponseWriter) error
}
