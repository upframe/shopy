package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// CartGET handles the GET request for /cart page
func CartGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return Redirect(w, r, "/login")
	}

	return RenderHTML(w, s, nil, "cart")
}
