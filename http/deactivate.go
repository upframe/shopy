package http

import (
	"database/sql"
	"net/http"
	"net/mail"
	"time"

	"github.com/upframe/fest"
	"github.com/upframe/fest/email"
)

// DeactivateHandler ...
type DeactivateHandler struct {
	SessionService fest.SessionService
	UserService    fest.UserService
	LinkService    fest.LinkService
}

func (h *DeactivateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (h *DeactivateHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	// Checks if the hash is indicated in the URL
	if r.URL.Query().Get("hash") == "" {
		return http.StatusNotFound, nil
	}

	// Fetches the link from the database
	link, err := h.LinkService.GetByHash(r.URL.Query().Get("hash"))

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
	err = h.UserService.Delete(link.User)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Marks the link as used and checks the errors
	link.Used = true
	err = h.LinkService.Update(link, "Used")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)
	return 0, nil
}

// POST ...
func (h *DeactivateHandler) POST(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := h.SessionService.Session(w, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !s.IsLoggedIn() {
		return http.StatusBadRequest, fest.ErrNotLoggedIn
	}

	// Sets the current time and expiration time of the deactivation email
	now := time.Now()
	expires := time.Now().Add(time.Hour * 2)

	link := &fest.Link{
		Path:    "/settings/deactivate",
		Hash:    fest.UniqueHash(s.Values["Email"].(string)),
		User:    s.Values["UserID"].(int),
		Used:    false,
		Time:    &now,
		Expires: &expires,
	}

	err = h.LinkService.Create(link)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	data := make(map[string]interface{})
	data["Name"] = s.User.FirstName + " " + s.User.LastName
	data["Hash"] = link.Hash
	data["Host"] = fest.BaseAddress

	email := &email.Email{
		From: &mail.Address{
			Name:    "Upframe",
			Address: email.FromDefaultEmail,
		},
		To: &mail.Address{
			Name:    data["Name"].(string),
			Address: s.User.Email,
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
