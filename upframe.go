package upframe

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/upframe/fest/pages"
)

// Upframe is the startup struct
type Upframe struct {
	Next httpserver.Handler
	Root string
}

func (u Upframe) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	// Checks if a static file (not directory) exists for this path. If it doesn't, we
	// handle the request.
	if info, err := os.Stat(u.Root + r.URL.Path); !(os.IsNotExist(err) || info.IsDir()) {
		return u.Next.ServeHTTP(w, r)
	}

	// Gets the current session or creates a new one if there is some error
	// decrypting it or if it doesn't exist
	s, _ := store.Get(r, "upframe-auth")

	// If it is a new session, initialize it, setting 'IsLoggedIn' as false
	if s.IsNew {
		s.Values["IsLoggedIn"] = false
	}

	// Saves the session in the cookie and checks for errors. This is useful
	// to reset the expiration time.
	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Routes the pages to the respective functions
	switch {
	case r.URL.Path == "/" && r.Method == http.MethodGet:
		return pages.IndexGET(w, r, s)
	case r.URL.Path == "/register" && r.Method == http.MethodGet:
		return pages.RegisterGET(w, r, s)
	case r.URL.Path == "/register" && r.Method == http.MethodPost:
		return pages.RegisterPOST(w, r, s)
	case r.URL.Path == "/login" && r.Method == http.MethodGet:
		return pages.LoginGET(w, r, s)
	case r.URL.Path == "/login" && r.Method == http.MethodPost:
		return pages.LoginPOST(w, r, s)
	case r.URL.Path == "/settings" && r.Method == http.MethodGet:
		return pages.SettingsGET(w, r, s)
	case r.URL.Path == "/store" && r.Method == http.MethodGet:
		return pages.StoreGET(w, r, s)
	case r.URL.Path == "/cart" && r.Method == http.MethodGet:
		return pages.CartGET(w, r, s)
	case r.URL.Path == "/checkout" && r.Method == http.MethodGet:
		return pages.CheckoutGET(w, r, s)
	case r.URL.Path == "/logout":
		return logout(w, r, s)
	}

	// Admin router: if the user is an admin and the page starts with /admin
	if pages.IsAdmin(s) && strings.HasPrefix(r.URL.Path, "/admin") {
		if strings.HasPrefix(r.URL.Path, "/admin/promocodes") {
			switch r.Method {
			case http.MethodGet:
				return pages.AdminPromocodesGET(w, r, s)
			case http.MethodPost:
				return pages.AdminPromocodesPOST(w, r)
			case http.MethodDelete:
				return pages.AdminPromocodesDELETE(w, r)
			case http.MethodPut:
				return pages.AdminPromocodesPUT(w, r)
			}
		}
	}

	// If the request doesn't match any route and it isn't a GET request
	// return a Status Not Implemented
	if r.Method != http.MethodGet {
		return http.StatusNotImplemented, nil
	}

	// Checks if there is a static template for this page. If so, show it!
	if _, err := os.Stat(filepath.Clean("templates/static" + r.URL.Path + ".tmpl")); err == nil {
		return pages.RenderHTML(w, nil, r.URL.Path)
	}

	// Return 404 Not Found for the rest
	return http.StatusNotFound, nil
}

// logout resets the session values and saves the cookie
func logout(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	// Reset the session values
	s.Values = map[interface{}]interface{}{}
	s.Values["IsLoggedIn"] = false

	// Saves the session and checks for error
	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return http.StatusOK, nil
}
