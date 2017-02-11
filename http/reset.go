package http

import (
	"database/sql"
	"net/http"
	"net/mail"
	"time"

	"github.com/bruhs/shopy"
)

// ResetGet ...
func ResetGet(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	s := r.Context().Value("session").(*shopy.Session)

	if hash := r.URL.Query().Get("hash"); hash != "" {
		// Fetches the link from the database
		link, err := c.Services.Link.Get(hash)

		// If the error is no rows, or the link is used, or it's expired or the path
		// is incorrect, show a 404 Not Found page.
		if err == sql.ErrNoRows || link.Used || link.Expires.Unix() < time.Now().Unix() || link.Path != "/reset" {
			return http.StatusNotFound, nil
		}

		// If there is any other error, return a 500
		if err != nil {
			return http.StatusInternalServerError, err
		}

		return Render(w, c, s, link.User, "reset/form")
	}

	return Render(w, c, s, nil, "reset/email")
}

// ResetPost ...
func ResetPost(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if hash := r.URL.Query().Get("hash"); hash != "" {
		// Fetches the link from the database
		var link *shopy.Link
		link, err = c.Services.Link.Get(hash)

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
		var user *shopy.User
		user, err = c.Services.User.Get(link.User)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		err = user.SetPassword(newPassword)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		err = c.Services.User.Update(user, "PasswordHash", "PasswordSalt")
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// SET LINK TO USED
		link.Used = true
		err = c.Services.Link.Update(link, "Used")
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
	var user *shopy.User
	user, err = c.Services.User.GetByEmail(formEmail)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Sets the current time and expiration time of the deactivation email
	now := time.Now()
	expires := time.Now().Add(time.Hour * 1)

	link := &shopy.Link{
		Path:    "/reset",
		Hash:    shopy.UniqueHash(formEmail),
		User:    user.ID,
		Used:    false,
		Time:    &now,
		Expires: &expires,
	}

	err = c.Services.Link.Create(link)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	data := make(map[string]interface{})
	data["Name"] = user.FirstName + " " + user.LastName
	data["Hash"] = link.Hash
	data["Host"] = c.BaseAddress

	email := &shopy.Email{
		From: &mail.Address{
			Name: "Upframe",
		},
		To: &mail.Address{
			Name:    data["Name"].(string),
			Address: formEmail,
		},
		Subject: "Reset your account password",
	}

	err = c.Services.Email.UseTemplate(email, data, "reset")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = c.Services.Email.Send(email)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
