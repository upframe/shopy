package pages

import (
	"net/http"

	"github.com/hacdias/upframe/models"
)

// RegisterGET handles the GET request for users
func RegisterGET(w http.ResponseWriter, r *http.Request) (int, error) {
	// TODO: Check url for referrer hash or show a "Registrations only available by invite" page

	return RenderHTML(w, nil, "register")
}

// RegisterPOST handles the POST http request in register page
func RegisterPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	// Check url for referrer hash and see if it's valid
	// otherwise return a http.StatusUnauthorized, nil

	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user := &models.User{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
		Invites:   0,
		Credit:    0,
		Confirmed: false,
	}

	if user.FirstName == "" || user.LastName == "" || user.Email == "" || r.FormValue("password") == "" {
		return http.StatusBadRequest, nil
	}

	if is, _ := isExistentUser(user.Email); is {
		return http.StatusConflict, nil
	}

	// TODO: Set referrer user

	user.GenerateReferrerHash()

	err = user.SetPassword(r.FormValue("password"))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = user.Insert()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// TODO: Send confirmation email
	return http.StatusCreated, nil
}

func isExistentUser(email string) (bool, error) {
	user, err := models.GetUserByEmail(email)
	if err != nil {
		return false, err
	}

	if user.ID != 0 {
		return true, nil
	}

	return false, nil
}
