package http

import (
	"log"
	"net/http"

	"github.com/upframe/fest"
)

func checkErrors(w http.ResponseWriter, r *http.Request, code int, err error) {
	if code != 0 {
		w.WriteHeader(code)
	}

	if err != nil {
		log.Print(err)
	}

	if r.Method == http.MethodGet {
		// TODO:
	}
}

// Redirect redirects the user to a page
func Redirect(w http.ResponseWriter, r *http.Request, path string) (int, error) {
	http.Redirect(w, r, path, http.StatusTemporaryRedirect)
	return http.StatusOK, nil
}

// GetSession ...
func GetSession(w http.ResponseWriter, r *http.Request, us fest.UserService) (*fest.Session, error) {
	// Create the session
	s := &fest.Session{}

	// Gets the current session or creates a new one if there is some error
	// decrypting it or if it doesn't exist
	s.Session, _ = fest.Store.Get(r, "upframe-auth")

	// If it is a new session, initialize it, setting 'IsLoggedIn' as false
	if s.IsNew {
		s.Values["IsLoggedIn"] = false
	}

	// Get the user info from the database and add it to the session data
	if s.IsLoggedIn() {
		var err error
		s.User, err = us.Get(s.Values["UserID"].(int))
		if err != nil {
			return s, err
		}
	}

	// Saves the session in the cookie and checks for errors. This is useful
	// to reset the expiration time.
	err := s.Save(r, w)
	if err != nil {
		return s, err
	}

	return s, nil
}
