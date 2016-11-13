package http

import (
	"net/http"

	"github.com/upframe/fest"
)

func checkErrors(w http.ResponseWriter, code int, err error) {

}

// Redirect redirects the user to a page
func Redirect(w http.ResponseWriter, r *http.Request, path string) (int, error) {
	http.Redirect(w, r, path, http.StatusTemporaryRedirect)
	return http.StatusOK, nil
}

func GetSession(w http.ResponseWriter, r *http.Request) (*fest.Session, error) {
	// Create the session
	s := &fest.Session{}

	// Gets the current session or creates a new one if there is some error
	// decrypting it or if it doesn't exist
	s.Session, _ = store.Get(r, "upframe-auth")

	// If it is a new session, initialize it, setting 'IsLoggedIn' as false
	if s.IsNew {
		s.Values["IsLoggedIn"] = false
	}

	// Get the user info from the database and add it to the session data
	if s.IsLoggedIn() {
		generic, err := models.GetUserByID(s.Values["UserID"].(int))
		if err != nil {
			return http.StatusInternalServerError, err
		}

		user := generic.(*models.User)
		s.User = user
	}

	// Saves the session in the cookie and checks for errors. This is useful
	// to reset the expiration time.
	err := s.Save(r, w)
	if err != nil {
		return s, err
	}

	return s, nil
}
