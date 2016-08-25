package upframe

import (
	"net/http"
	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// Upframe is the startup struct
type Upframe struct {
	Next httpserver.Handler
}

func (u Upframe) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if shouldHandle(r.URL.Path) {
		w.Write([]byte("We are Upframe"))

		return http.StatusOK, nil
	}

	return u.Next.ServeHTTP(w, r)
}

func shouldHandle(path string) bool {
	paths := []string{"/css", "/js", "/imgs"}

	for i := range paths {
		if strings.HasPrefix(path, paths[i]) {
			return false
		}
	}

	return true
}
