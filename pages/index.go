package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// IndexGET handles the GET request for /index page
func IndexGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	return RenderHTML(w, s, nil, "index")
}
