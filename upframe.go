package upframe

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/hacdias/upframe/pages"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// Upframe is the startup struct
type Upframe struct {
	Next httpserver.Handler
}

func (u Upframe) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	// Checks if a static file (not directory) exists for this path. If it doesn't, we
	// handle the request.
	if info, err := os.Stat("static" + r.URL.Path); os.IsNotExist(err) || info.IsDir() {
		switch r.URL.Path {
		case "/":
			return pages.IndexGET(w, r)
		case "/register":
			// if logged in redirect to / or /store
			// else:

			if r.Method == http.MethodGet {
				return pages.RegisterGET(w, r)
			} else if r.Method == http.MethodPost {
				return pages.RegisterPOST(w, r)
			}

			return http.StatusNotImplemented, nil
		case "/login":
			// if logged in redirect to / or /store

			if r.Method == http.MethodGet {
				return pages.LoginGET(w, r)
			} else if r.Method == http.MethodPost {
				return pages.LoginPOST(w, r)
			}

			return http.StatusNotImplemented, nil
		case "/settings":
			// if not logged in redirect to /login
			return pages.SettingsGET(w, r)
		case "/store":
			//return utils.RenderHTML(w, nil, "store")
			return pages.StoreGET(w, r)
		case "/cart":
			// if not logged in redirect to /login
			return pages.CartGET(w, r)
		case "/checkout":
			/// if not logged in redirect to /login
			return pages.CheckoutGET(w, r)
		}

		// Checks if there is a static template for this page. If so, show it!
		if _, err := os.Stat(filepath.Clean("templates/static" + r.URL.Path + ".tmpl")); err == nil {
			return pages.RenderHTML(w, nil, r.URL.Path)
		}

		return http.StatusNotFound, nil
	}

	return u.Next.ServeHTTP(w, r)
}
