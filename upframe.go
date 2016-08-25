package upframe

import (
	"net/http"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type Upframe struct {
	Next httpserver.Handler
}

func (u Upframe) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {

	return u.Next.ServeHTTP(w, r)
}
