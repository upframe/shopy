package pages

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/hacdias/upframe/models"
)

// RegisterGET handles the GET request for register page
func RegisterGET(w http.ResponseWriter, r *http.Request) (int, error) {
	// If we're in the success page, show it, LOL
	if r.URL.Query().Get("success") != "" {
		return RenderHTML(w, nil, "register-success")
	}

	// Gets the referrer user
	referrer, err := models.GetUserByReferral(r.URL.Query().Get("ref"))

	// If the user doesn't exist show a page telling that registration
	// is invitation only
	if err != nil {
		return RenderHTML(w, nil, "register-invite")
	}

	// If the user exists, but doesn't have invites, show that information
	if referrer.Invites < 1 {
		return RenderHTML(w, referrer, "register-gone")
	}

	// Otherwise, show the registration page
	return RenderHTML(w, nil, "register")
}

// RegisterPOST handles the POST http request in register page
func RegisterPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	// Gets the referrer user using the ?referral= option in the URL. If it doesn't
	// find the user, return a 403 Forbidden status
	referrer, err := models.GetUserByReferral(r.URL.Query().Get("ref"))
	if err != nil {
		return http.StatusForbidden, nil
	}

	// Checks if the referrer still has invites available! This is important! If
	// it doesn't, return a 410 Status Gone
	if referrer.Invites < 1 {
		return http.StatusGone, nil
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

	// Decrement one value from the referrer invites number and updates it in
	// the database and checks for errors. In this case, if there is an error
	// we will keep the registration going because it's no user-fault and the
	// refferer had invites available.
	//
	// The error is logged using a prefix and can be checked afterwards by
	// system administrators.
	referrer.Invites--
	err = referrer.Update("invites")

	if err != nil {
		log.Println("INVITE DECREMENT ERROR: " + err.Error())
	}
	/*
		// TODO: Send confirmation email
		now := time.Now()
		expires := now.Add(models.ConfirmExpiration)

		link := &models.Link{
			Path:    "/register",
			Hash:    models.UniqueHash(u.Email),
			User:    u.ID,
			Used:    false,
			Time:    &now,
			Expires: &expires,
		}

		err := link.Insert()
		if err != nil {
			return err
		}

		data := make(map[string]interface{})
		data["Name"] = user.FirstName + " " + user.LastName
		data["Hash"] = link.Hash
		data["Host"] = r.URL.Scheme + "://" + r.URL.Host */

	// TODO: Finish

	// Redirect to success page as a Temporary Redirect
	// TODO: JavaScript "popup" to show this information instead of redirect?
	http.Redirect(w, r, "/register?success=true", http.StatusTemporaryRedirect)
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
