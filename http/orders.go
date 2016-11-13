package http

import (
	"net/http"

	"github.com/upframe/fest"
)

// OrdersHandler ...
type OrdersHandler struct {
	OrderService fest.OrderService
	UserService  fest.UserService
}

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
	s, err := GetSession(w, r, h.UserService)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !s.IsLoggedIn() {
		return Redirect(w, r, "/login")
	}

	data, err := h.OrderService.GetByUser(s.Values["UserID"].(int))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, data, "orders")
}
