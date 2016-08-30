package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// SettingsGET handles the GET request for /settings page
func SettingsGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !isLoggedIn(s) {
		return redirect(w, r, "/login")
	}

	return RenderHTML(w, nil, "settings")
}
