package http

import (
	"net/http"

	"github.com/upframe/fest"
)

// IndexHandler ...
type IndexHandler handler

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)
	defer checkErrors(w, r, code, err)

	switch r.Method {
	case http.MethodGet:
		code, err = h.GET(w, r)
	default:
		code, err = http.StatusNotImplemented, nil
	}
}

// GET ...
func (h *IndexHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	return RenderHTML(w, s, nil, "index")
}
