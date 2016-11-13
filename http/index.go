package http

import (
	"net/http"

	"github.com/upframe/fest"
)

// IndexHandler ...
type IndexHandler struct {
	UserService fest.UserService
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)

	switch r.Method {
	case http.MethodGet:
		code, err = h.GET(w, r)
	default:
		code, err = http.StatusNotImplemented, nil
	}

	checkErrors(w, code, err)
}

// GET ...
func (h *IndexHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := GetSession(w, r, h.UserService)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, nil, "index")
}
