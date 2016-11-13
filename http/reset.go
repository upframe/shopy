package http

import (
	"database/sql"
	"net/http"
	"net/mail"
	"time"

	"github.com/upframe/fest"
	"github.com/upframe/fest/email"
)

// ResetHandler ...
type ResetHandler struct {
	UserService fest.UserService
	LinkService fest.LinkService
}

func (h *ResetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)

	switch r.Method {
	case http.MethodGet:
		code, err = h.GET(w, r)
	case http.MethodPost:
		code, err = h.POST(w, r)
	default:
		code, err = http.StatusNotImplemented, nil
	}

	checkErrors(w, code, err)
}

// GET ...
func (h *ResetHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := GetSession(w, r, h.UserService)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if hash := r.URL.Query().Get("hash"); hash != "" {
		// Fetches the link from the database
		link, err := h.LinkService.Get(hash)

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

// POST ...
func (h *ResetHandler) POST(w http.ResponseWriter, r *http.Request) (int, error) {
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if hash := r.URL.Query().Get("hash"); hash != "" {
		// Fetches the link from the database
		var link *fest.Link
		link, err = h.LinkService.Get(hash)

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
		var user *fest.User
		user, err = h.UserService.Get(link.User)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		err = user.SetPassword(newPassword)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		err = h.UserService.Update(user, "PasswordHash", "PasswordSalt")
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// SET LINK TO USED
		link.Used = true
		err = h.LinkService.Update(link, "Used")
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
	var user *fest.User
	user, err = h.UserService.GetByEmail(formEmail)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Sets the current time and expiration time of the deactivation email
	now := time.Now()
	expires := time.Now().Add(time.Hour * 1)

	link := &fest.Link{
		Path:    "/reset",
		Hash:    fest.UniqueHash(formEmail),
		User:    user.ID,
		Used:    false,
		Time:    &now,
		Expires: &expires,
	}

	err = h.LinkService.Create(link)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	data := make(map[string]interface{})
	data["Name"] = user.FirstName + " " + user.LastName
	data["Hash"] = link.Hash
	data["Host"] = fest.BaseAddress

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

	return http.StatusOK, nil
}
