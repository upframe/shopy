package http

import (
	"context"
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

// InjectSession ...
func InjectSession(h http.Handler, us fest.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				checkErrors(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		// Saves the session in the cookie and checks for errors. This is useful
		// to reset the expiration time.
		err := s.Save(r, w)
		if err != nil {
			checkErrors(w, r, http.StatusInternalServerError, err)
			return
		}

		ctx := context.WithValue(r.Context(), "session", s)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}
