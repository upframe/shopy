package pages

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/upframe/fest/models"
)

var (
	errAlreadyLoggedIn = errors.New("The user is already logged in.")
	errNotLoggedIn     = errors.New("The user is not logged in.")
)

// LoginGET handles the GET request for /login page
func LoginGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if IsLoggedIn(s) {
		return Redirect(w, r, "/")
	}

	return RenderHTML(w, s, nil, "login")
}

// LoginPOST handles the POST request for /login page
func LoginPOST(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if IsLoggedIn(s) {
		return http.StatusBadRequest, errAlreadyLoggedIn
	}

	if r.Header.Get("Resend") == "true" {
		// Obtains the user and checks for errors
		user, err := models.GetUserByEmail(r.Header.Get("Email"))
		if err == sql.ErrNoRows {
			return http.StatusNotFound, err
		}

		if err != nil {
			return http.StatusInternalServerError, err
		}

		return confirmationEmail(user)
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
	user, err := models.GetUserByEmail(email)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
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

	// Checks if the user is deactivated
	if user.Deactivated {
		return http.StatusLocked, nil
	}

	// Sets the session cookie values
	s.Values["IsLoggedIn"] = true
	s.Values["IsAdmin"] = user.Admin
	s.Values["UserID"] = user.ID
	s.Values["FirstName"] = user.FirstName
	s.Values["LastName"] = user.LastName
	s.Values["Email"] = user.Email

	// Initialize cart
	s.Values["Cart"] = &cart{
		Products: map[int]*cartItem{},
	}

	// Saves the cookie and checks for errors
	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
