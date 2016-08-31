package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// SettingsGET handles the GET request for /settings page
func SettingsGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return Redirect(w, r, "/login")
	}

	return RenderHTML(w, s, nil, "settings")
}
