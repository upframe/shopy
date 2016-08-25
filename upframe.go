package upframe

import (
	"net/http"
	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type Upframe struct {
	Next httpserver.Handler
}

func (u Upframe) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if !strings.HasPrefix(r.URL.Path, "/assets") {
		w.Write([]byte("We are Upframe"))

		return http.StatusOK, nil
	}

	return u.Next.ServeHTTP(w, r)
}
