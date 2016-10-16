package pages

import (
	"net/http"

	"github.com/upframe/fest/models"
)

// AdminPromocodesGET handles the GET request for every /admin/promocodes/... URLs
func AdminPromocodesGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	return AdminGenericGET(w, r, s, "promocodes", models.GetPromocodes)
}

// AdminPromocodesPOST creates a new item
func AdminPromocodesPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericPOST(w, r, new(models.Promocode))
}

// AdminPromocodesDELETE deactivates a promocode
func AdminPromocodesDELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericDELETE(w, r, "promocodes", models.GetPromocode)
}

// AdminPromocodesPUT changes a promocode
func AdminPromocodesPUT(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericPUT(w, r, new(models.Promocode), models.UpdateAll)
}
