package http

import (
	"net/http"

	"github.com/upframe/fest"
)

// IndexGet ...
func IndexGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	return RenderHTML(w, c, s, nil, "index")
}
