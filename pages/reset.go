package pages

import (
	"database/sql"
	"net/http"
	"net/mail"
	"time"

	"github.com/upframe/fest/email"
	"github.com/upframe/fest/models"
)

// ResetGET handles GET request to display the reset password page
func ResetGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	if hash := r.URL.Query().Get("hash"); hash != "" {
		// Fetches the link from the database
		link, err := models.GetLinkByHash(hash)

		// If the error is no rows, or the link is used, or it's expired or the path
		// is incorrect, show a 404 Not Found page.
		if err == sql.ErrNoRows || link.Used || link.Expires.Unix() < time.Now().Unix() || link.Path != "/reset" {
			return http.StatusNotFound, nil
		}

		// If there is any other error, return a 500
		if err != nil {
			return http.StatusInternalServerError, err
		}

		return RenderHTML(w, s, link.User, "reset/form")
	}
	return RenderHTML(w, s, nil, "reset/email")
}

// ResetPOST sends an email to the user Email
func ResetPOST(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if hash := r.URL.Query().Get("hash"); hash != "" {
		// Fetches the link from the database
		var link *models.Link
		link, err = models.GetLinkByHash(hash)

		// If the error is no rows, or the link is used, or it's expired or the path
		// is incorrect, show a 404 Not Found page.
		if err == sql.ErrNoRows || link.Used || link.Expires.Unix() < time.Now().Unix() || link.Path != "/reset" {
			return http.StatusNotFound, nil
		}

		// If there is any other error, return a 500
		if err != nil {
			return http.StatusInternalServerError, err
		}

		newPassword := r.FormValue("password")
		if newPassword == "" {
			return http.StatusBadRequest, nil
		}

		// SET USER PASSWORD AND UPDATE PWD HASH AND PWD SALT
		var g models.Generic
		g, err = models.GetUserByID(link.User)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		user := g.(*models.User)
		err = user.SetPassword(newPassword)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		err = user.Update("password_hash", "password_salt")
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// SET LINK TO USED
		link.Used = true
		err = link.Update("used")
		if err != nil {
			return http.StatusInternalServerError, err
		}

		return http.StatusOK, nil
	}

	// get email from form
	formEmail := r.FormValue("email")

	if formEmail == "" {
		return http.StatusBadRequest, nil
	}

	// get user by email
	user, err := models.GetUserByEmail(formEmail)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Sets the current time and expiration time of the deactivation email
	now := time.Now()
	expires := time.Now().Add(time.Hour * 1)

	link := &models.Link{
		Path:    "/reset",
		Hash:    models.UniqueHash(formEmail),
		User:    user.ID,
		Used:    false,
		Time:    &now,
		Expires: &expires,
	}

	err = link.Insert()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	data := make(map[string]interface{})
	data["Name"] = user.FirstName + " " + user.LastName
	data["Hash"] = link.Hash
	data["Host"] = BaseAddress

	email := &email.Email{
		From: &mail.Address{
			Name:    "Upframe",
			Address: email.FromDefaultEmail,
		},
		To: &mail.Address{
			Name:    data["Name"].(string),
			Address: formEmail,
		},
		Subject: "Reset your account password",
	}

	err = email.UseTemplate("reset", data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = email.Send()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}
