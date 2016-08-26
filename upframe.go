package upframe

import (
	"net/http"
	"os"

	"github.com/hacdias/upframe/utils"
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
			return utils.RenderHTML(w, nil, "index")
		case "/register":
			// if logged in redirect to / or /store
			return utils.RenderHTML(w, nil, "register")
		case "/login":
			// if logged in redirect to / or /store
			return utils.RenderHTML(w, nil, "login")
		case "/settings":
			// if not logged in redirect to /login
			return utils.RenderHTML(w, nil, "settings")
		case "/store":
			return utils.RenderHTML(w, nil, "store")
		case "/cart":
			// if not logged in redirect to /login
			return utils.RenderHTML(w, nil, "cart")
		case "/checkout":
			/// if not logged in redirect to /login
			return utils.RenderHTML(w, nil, "checkout")
		case "/wishlist":
			// if not logged in redirect to /login
			return utils.RenderHTML(w, nil, "wishlist")
		}

		return http.StatusNotFound, nil
	}

	return u.Next.ServeHTTP(w, r)
}
