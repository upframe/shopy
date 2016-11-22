package http

import (
	"database/sql"
	"net/http"
	"net/mail"
	"time"

	"github.com/upframe/fest"
)

// DeactivateGet ...
func DeactivateGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	// Checks if the hash is indicated in the URL
	if r.URL.Query().Get("hash") == "" {
		return http.StatusNotFound, nil
	}

	// Fetches the link from the database
	link, err := c.Services.Link.Get(r.URL.Query().Get("hash"))

	// If the error is no rows, or the link is used, or it's expired or the path
	// is incorrect, show a 404 Not Found page.
	if err == sql.ErrNoRows || link.Used || link.Expires.Unix() < time.Now().Unix() || link.Path != "/settings/deactivate" {
		return http.StatusNotFound, nil
	}

	// If there is any other error, return a 500
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Deactivates the user and checks for error
	err = c.Services.User.Delete(link.User)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Marks the link as used and checks the errors
	link.Used = true
	err = c.Services.Link.Update(link, "Used")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)
	return 0, nil
}

// DeactivatePost ...
func DeactivatePost(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	if !s.Logged {
		return http.StatusBadRequest, fest.ErrNotLoggedIn
	}

	// Sets the current time and expiration time of the deactivation email
	now := time.Now()
	expires := time.Now().Add(time.Hour * 2)

	link := &fest.Link{
		Path:    "/settings/deactivate",
		Hash:    fest.UniqueHash(s.User.Email),
		User:    s.User.ID,
		Used:    false,
		Time:    &now,
		Expires: &expires,
	}

	err := c.Services.Link.Create(link)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	data := make(map[string]interface{})
	data["Name"] = s.User.FirstName + " " + s.User.LastName
	data["Hash"] = link.Hash
	data["Host"] = c.BaseAddress

	email := &fest.Email{
		From: &mail.Address{
			Name: "Upframe",
		},
		To: &mail.Address{
			Name:    data["Name"].(string),
			Address: s.User.Email,
		},
		Subject: "Deactivate your account",
	}

	err = c.Services.Email.UseTemplate(email, data, "deactivation")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = c.Services.Email.Send(email)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
