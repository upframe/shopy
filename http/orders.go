package http

import (
	"net/http"
	"strconv"

	"github.com/upframe/fest"
)

// OrdersGet ...
func OrdersGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	if !s.IsLoggedIn() {
		return Redirect(w, r, "/login")
	}

	data, err := c.Services.Order.GetsWhere(0, 0, "ID", "UserID", strconv.Itoa(s.Values["UserID"].(int)))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Render(w, c, s, data, "orders")
}
