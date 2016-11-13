package http

import (
	"net/http"

	"github.com/upframe/fest"
)

// StoreHandler ...
type StoreHandler struct {
	ProductsService fest.ProductService
	SessionService  fest.SessionService
}

func (h *StoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (h *StoreHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := h.SessionService.Session(w, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	products, err := h.ProductsService.GetsWhere(0, 0, "name", "deactivated", "0")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, products, "store")
}
