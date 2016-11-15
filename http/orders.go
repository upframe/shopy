package http

import (
	"net/http"

	"github.com/upframe/fest"
)

// OrdersGet ...
func OrdersGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	if !s.IsLoggedIn() {
		return Redirect(w, r, "/login")
	}

	data, err := c.Services.Order.GetByUser(s.Values["UserID"].(int))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, c, s, data, "orders")
}
