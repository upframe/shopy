package pages

import (
	"database/sql"
	"net/http"

	"github.com/hacdias/upframe/models"
)

// RegisterGET handles the GET request for register page
func RegisterGET(w http.ResponseWriter, r *http.Request) (int, error) {
	// TODO: Check url for referrer hash or show a "Registrations only available by invite" page

	return RenderHTML(w, nil, "register")
}

// RegisterPOST handles the POST http request in register page
func RegisterPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	// Gets the referrer user using the ?referral= option in the URL. If it doesn't
	// find the user, return a 403 Forbidden status.
	referrer, err := models.GetUserByReferral(r.URL.Query().Get("ref"))
	if err != nil {
		return http.StatusForbidden, nil
	}

	// Parses the form and checks for errors
	err = r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Builds the user variable
	user := &models.User{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
		Invites:   0,
		Credit:    0,
		Confirmed: false,
		Referrer:  sql.NullInt64{Int64: int64(referrer.ID), Valid: true},
	}

	// Checks if any of the fields is empty, if so, return a 400 Bad Request error
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || r.FormValue("password") == "" {
		return http.StatusBadRequest, nil
	}

	// Checks if there is already an user with this email. If there is,
	// return a 407 Conflict error.
	if is, _ := isExistentUser(user.Email); is {
		return http.StatusConflict, nil
	}

	// Generates a unique referral hash for this user
	user.GenerateReferralHash()

	// Sets the password hash and salt for the user and checks for errors
	err = user.SetPassword(r.FormValue("password"))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Inserts the user into the database
	err = user.Insert()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// TODO: Send confirmation email

	return http.StatusCreated, nil
}

// isExistentUser checks if there is an user with the specified email
// and returns true and nil if the user exists and there is no error
func isExistentUser(email string) (bool, error) {
	// Fetches the user from the database and checks for errors
	user, err := models.GetUserByEmail(email)
	if err != nil {
		return false, err
	}

	// Checks if the user ID is different from 0, which means that it is valid
	// if so, returns true and nil
	return (user.ID != 0), nil
}
