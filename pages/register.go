package pages

import "net/http"

func RegisterGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "register")
}
