package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/upframe/fest/models"
)

// StoreGET handles the GET request for /store page
func StoreGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	products, err := models.GetProductsWhere(0, 0, "name", "deactivated", "0")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, products, "store")
}
