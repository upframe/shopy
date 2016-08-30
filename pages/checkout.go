package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// CheckoutGET handles the GET request for /checkout page
func CheckoutGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !isLoggedIn(s) {
		return redirect(w, r, "/login")
	}

	return RenderHTML(w, nil, "checkout")
}
