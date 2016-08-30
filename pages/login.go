package pages

import (
	"database/sql"
	"net/http"

	"github.com/hacdias/upframe/models"
)

// LoginGET handles the GET request for /login page
func LoginGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "login")
}

// LoginPOST handles the POST request for /login page
func LoginPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	// Parses the form and checks for errors
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		return http.StatusBadRequest, nil
	}

	user, err := models.GetUserByEmail(email)

	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	ok, err := user.CheckPassword(password)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !ok {
		return http.StatusUnauthorized, nil
	}

	return http.StatusOK, nil
}
