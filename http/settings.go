package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/upframe/fest"
)

// SettingsHandler ...
type SettingsHandler struct {
	SessionService fest.SessionService
	UserService    fest.UserService
}

func (h *SettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

type settings struct {
	User    *fest.User
	BaseURL string
}

// GET ...
func (h *SettingsHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := h.SessionService.Session(w, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !s.IsLoggedIn() {
		return Redirect(w, r, "/login")
	}

	user, err := h.UserService.Get(s.Values["UserID"].(int))
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, settings{
		User:    user,
		BaseURL: BaseAddress,
	}, "settings")
}

// POST ...
func (h *SettingsHandler) POST(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := h.SessionService.Session(w, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !s.IsLoggedIn() {
		return http.StatusBadRequest, fest.ErrNotLoggedIn
	}

	// Get the JSON information
	rawBuffer := &bytes.Buffer{}
	rawBuffer.ReadFrom(r.Body)

	// Parses the JSON into a user object and checks for errors
	user := &fest.User{}
	err = json.Unmarshal(rawBuffer.Bytes(), user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if user.ID == 0 {
		return http.StatusBadRequest, errors.New("The ID of the object isn't set.")
	}

	// Updates and checks for errors
	err = h.UserService.Update(user, "first_name", "last_name", "email", "address")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
