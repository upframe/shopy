package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/upframe/fest"
)

type message struct {
	Code    int
	Message string
	Error   error `json:"-"`
}

// FestHandler ...
type FestHandler func(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error)

// Inject ...
func Inject(h FestHandler, c *fest.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			code int
			err  error
		)

		defer func() {
			if code == 0 && err == nil {
				return
			}

			msg := &message{Code: code}

			if err != nil {
				c.Logger.Print(err)
				msg.Message = err.Error()
			} else {
				msg.Message = http.StatusText(code)
			}

			if code != 0 {
				w.WriteHeader(code)
			}

			if strings.HasPrefix(r.URL.Path, "/api") || r.Method != http.MethodGet {
				data, e := json.MarshalIndent(msg, "", "\t")
				if e != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}

				w.Write(data)
				return
			}

			// TODO: show page
			w.Write([]byte(msg.Message))
		}()

		// Create the session
		s := &fest.Session{}

		// Gets the current session or creates a new one if there is some error
		// decrypting it or if it doesn't exist
		s.Session, _ = c.Store.Get(r, "upframe-auth")

		// If it is a new session, initialize it, setting 'IsLoggedIn' as false
		if s.IsNew {
			s.Values["IsLoggedIn"] = false
		}

		// Get the user info from the database and add it to the session data
		if s.IsLoggedIn() {
			s.User, err = c.Services.User.Get(s.Values["UserID"].(int))
			if err != nil {
				return
			}
		}

		// Saves the session in the cookie and checks for errors. This is useful
		// to reset the expiration time.
		err = s.Save(r, w)
		if err != nil {
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", s))
		code, err = h(w, r, c)
	}
}

// MustLogin ...
func MustLogin(h FestHandler) FestHandler {
	return func(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
		s := r.Context().Value("session").(*fest.Session)

		if s.IsLoggedIn() {
			return h(w, r, c)
		}

		if r.Method == http.MethodGet && !strings.HasPrefix(r.URL.Path, "/api") {
			return Redirect(w, r, "/login")
		}

		return http.StatusUnauthorized, fest.ErrNotLoggedIn
	}
}

// MustAdmin ...
func MustAdmin(h FestHandler) FestHandler {
	return MustLogin(func(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
		s := r.Context().Value("session").(*fest.Session)

		if s.IsAdmin() {
			return h(w, r, c)
		}

		return http.StatusForbidden, nil
	})
}

// Redirect redirects the user to a page
func Redirect(w http.ResponseWriter, r *http.Request, path string) (int, error) {
	http.Redirect(w, r, path, http.StatusTemporaryRedirect)
	return 0, nil
}
