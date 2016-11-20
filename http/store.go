package http

import (
	"net/http"

	"github.com/upframe/fest"
)

// StoreGet ...
func StoreGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.SessionCookie)

	products, err := c.Services.Product.GetsWhere(0, 0, "Name", "Deactivated", "0")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Render(w, c, s, products, "store")
}
