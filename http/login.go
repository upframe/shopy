package http

import (
	"database/sql"
	"net/http"

	"github.com/upframe/fest"
)

// LoginGet ...
func LoginGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	if s.Logged {
		return Redirect(w, r, "/")
	}

	return Render(w, c, s, nil, "login")
}

// LoginPost ...
func LoginPost(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	if s.Logged {
		return http.StatusBadRequest, fest.ErrAlreadyLoggedIn
	}

	if r.Header.Get("Resend") == "true" {
		// Obtains the user and checks for errors
		user, err := c.Services.User.GetByEmail(r.Header.Get("email"))
		if err == sql.ErrNoRows {
			return http.StatusNotFound, err
		}

		if err != nil {
			return http.StatusInternalServerError, err
		}

		return confirmationEmail(c, user)
	}

	// Parses the form and checks for errors
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Obtains the email and the password
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Checks if they're blank or not
	if email == "" || password == "" {
		return http.StatusBadRequest, nil
	}

	// Obtains the user and checks for errors
	user, err := c.Services.User.GetByEmail(email)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Checks if the user is deactivated
	if user.Deactivated {
		return http.StatusLocked, nil
	}

	// Checks the password and checks for errors
	ok, err := user.CheckPassword(password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !ok {
		return http.StatusUnauthorized, nil
	}

	// Checks if the user is confirmed
	if !user.Confirmed {
		return http.StatusFailedDependency, nil
	}

	// Sets the session cookie values
	s.Logged = true
	s.User.ID = user.ID

	// Saves the cookie and checks for errors
	err = c.Services.Session.Save(w, s)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
