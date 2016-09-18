package pages

import (
	"database/sql"
	"net/http"
	"net/mail"
	"time"

	"github.com/gorilla/sessions"
	"github.com/upframe/fest/email"
	"github.com/upframe/fest/models"
)

// DeactivateGET creates a new deactivation link
func DeactivateGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	// Checks if the hash is indicated in the URL
	if r.URL.Query().Get("deactivate") == "" {
		return http.StatusNotFound, nil
	}

	// Fetches the link from the database
	link, err := models.GetLinkByHash(r.URL.Query().Get("hash"))

	// If the error is no rows, or the link is used, or it's expired or the path
	// is incorrect, show a 404 Not Found page.
	if err == sql.ErrNoRows || link.Used || link.Expires.Unix() < time.Now().Unix() || link.Path != "/settings/deactivate" {
		return http.StatusNotFound, nil
	}

	// If there is any other error, return a 500
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Gets the users and checks for error
	g, err := models.GetUserByID(link.User)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Deactivates the user and checks for error
	user := g.(*models.User)
	err = user.Deactivate()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Marks the link as used and checks the errors
	link.Used = true
	err = link.Update("used")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return http.StatusOK, nil
}

// DeactivatePOST creates the deactivation email and sends it to the user
func DeactivatePOST(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return http.StatusBadRequest, errNotLoggedIn
	}

	// Sets the current time and expiration time of the deactivation email
	now := time.Now()
	expires := time.Now().Add(time.Hour * 2)

	link := &models.Link{
		Path:    "/settings/deactivate",
		Hash:    models.UniqueHash(s.Values["Email"].(string)),
		User:    s.Values["UserID"].(int),
		Used:    false,
		Time:    &now,
		Expires: &expires,
	}

	err := link.Insert()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	data := make(map[string]interface{})
	data["Name"] = s.Values["FirstName"].(string) + " " + s.Values["LastName"].(string)
	data["Hash"] = link.Hash
	data["Host"] = BaseAddress

	email := &email.Email{
		From: &mail.Address{
			Name:    "Upframe",
			Address: email.FromDefaultEmail,
		},
		To: &mail.Address{
			Name:    data["Name"].(string),
			Address: s.Values["Email"].(string),
		},
		Subject: "Deactivate your account",
	}

	err = email.UseTemplate("deactivation", data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = email.Send()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
