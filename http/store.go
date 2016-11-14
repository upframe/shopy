package http

import (
	"net/http"

	"github.com/upframe/fest"
)

// StoreHandler ...
type StoreHandler handler

func (h *StoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (h *StoreHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	products, err := h.Services.Product.GetsWhere(0, 0, "name", "deactivated", "0")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, products, "store")
}
