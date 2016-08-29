package pages

import "net/http"

// LoginGET handles the GET request for /login page
func LoginGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "login")
}
