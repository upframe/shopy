package http

import (
	"net/http"

	"github.com/upframe/fest"
)

// OrdersHandler ...
type OrdersHandler handler

func (h *OrdersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (h *OrdersHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	if !s.IsLoggedIn() {
		return Redirect(w, r, "/login")
	}

	data, err := h.Services.Order.GetByUser(s.Values["UserID"].(int))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, data, "orders")
}
