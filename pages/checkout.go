package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// CheckoutGET handles the GET request for /checkout page
func CheckoutGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return Redirect(w, r, "/login")
	}

	return RenderHTML(w, nil, "checkout")
}
