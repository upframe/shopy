package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/upframe/fest"
)

type settings struct {
	User    *fest.User
	BaseURL string
}

// SettingsGet ...
func SettingsGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	user, err := c.Services.User.Get(s.Values["UserID"].(int))
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Render(w, c, s, settings{
		User:    user,
		BaseURL: c.BaseAddress,
	}, "settings")
}

// SettingsPost ...
func SettingsPost(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	// Get the JSON information
	rawBuffer := &bytes.Buffer{}
	rawBuffer.ReadFrom(r.Body)

	// Parses the JSON into a user object and checks for errors
	user := &fest.User{}
	err := json.Unmarshal(rawBuffer.Bytes(), user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if user.ID == 0 {
		return http.StatusBadRequest, errors.New("The ID of the object isn't set.")
	}

	// Updates and checks for errors
	err = c.Services.User.Update(user, "FirstName", "LastName", "Email", "Address")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
