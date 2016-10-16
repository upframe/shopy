package pages

import (
	"net/http"

	"github.com/upframe/fest/models"
)

// MyOrdersGET displays user orders
func MyOrdersGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	if !s.IsLoggedIn() {
		return Redirect(w, r, "/login")
	}

	data, err := models.GetAllOrdersByUser(s.Values["UserID"].(int))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, data, "myorder")
}
