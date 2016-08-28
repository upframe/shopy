package pages

import "net/http"

func SettingsGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "settings")
}
