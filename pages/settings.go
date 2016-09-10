package pages

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/upframe/fest/models"
)

// SettingsGET handles the GET request for /settings page
func SettingsGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return Redirect(w, r, "/login")
	}

	user, err := models.GetUserByID(s.Values["UserID"].(int))
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, user, "settings")
}

// SettingsPUT handles the PUT request for /settings page which is the method
// for updating the user information
func SettingsPUT(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return http.StatusBadRequest, errNotLoggedIn
	}

	// Get the JSON information
	rawBuffer := &bytes.Buffer{}
	rawBuffer.ReadFrom(r.Body)

	// Parses the JSON into a user object and checks for errors
	user := &models.User{}
	err := json.Unmarshal(rawBuffer.Bytes(), user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if user.ID == 0 {
		return http.StatusBadRequest, errors.New("The ID of the object isn't set.")
	}

	// Updates and checks for errors
	err = user.Update("first_name", "last_name", "email", "address")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
