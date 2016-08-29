package pages

import "net/http"

// SettingsGET handles the GET request for /settings page
func SettingsGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "settings")
}
