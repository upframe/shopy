package pages

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/hacdias/upframe/models"
)

// LoginGET handles the GET request for /login page
func LoginGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if isLoggedIn(s) {
		return redirect(w, r, "/")
	}

	return RenderHTML(w, nil, "login")
}

// LoginPOST handles the POST request for /login page
func LoginPOST(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if isLoggedIn(s) {
		return http.StatusBadRequest, nil
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

	// Sets the session cookie
	s.Values["logged"] = true
	s.Values["uid"] = user.ID
	s.Values["admin"] = user.Admin

	// Saves the cookie and checks for errors
	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
