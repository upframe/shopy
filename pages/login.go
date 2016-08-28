package pages

import "net/http"

func LoginGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "login")
}
