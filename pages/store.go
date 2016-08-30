package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// StoreGET handles the GET request for /store page
func StoreGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	return RenderHTML(w, nil, "store")
}
