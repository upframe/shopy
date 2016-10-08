package pages

import (
	"net/http"

	"github.com/upframe/fest/models"
)

const itemsPerPage = 50

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

// AdminOrdersGET is
func AdminOrdersGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	return AdminGenericGET(w, r, s, "orders", models.GetOrders)
}

// AdminOrdersPOST is
func AdminOrdersPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericPOST(w, r, new(models.Order))
}

// AdminOrdersDELETE is
func AdminOrdersDELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericDELETE(w, r, "orders", models.GetOrder)
}

// AdminOrdersPUT is
func AdminOrdersPUT(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericPUT(w, r, new(models.Order), models.UpdateAll)
}

// AdminUsersGET is
func AdminUsersGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	return AdminGenericGET(w, r, s, "users", models.GetUsers)
}

// AdminUsersPOST is
func AdminUsersPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericPOST(w, r, new(models.User))
}

// AdminUsersDELETE is
func AdminUsersDELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericDELETE(w, r, "users", models.GetUserByID)
}

// AdminUsersPUT is
func AdminUsersPUT(w http.ResponseWriter, r *http.Request) (int, error) {
	// In Users, we do not update: password_hash, password_salt, referral,
	// nor referrer
	return AdminGenericPUT(w, r, new(models.User), "id",
		"first_name",
		"last_name",
		"email",
		"address",
		"invites",
		"credit",
		"confirmed",
		"admin",
		"deactivated")
}

// AdminProductsGET is
func AdminProductsGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	return AdminGenericGET(w, r, s, "products", models.GetProducts)
}

// AdminProductsPOST is
func AdminProductsPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericPOST(w, r, new(models.Product))
}

// AdminProductsDELETE is
func AdminProductsDELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericDELETE(w, r, "products", models.GetProduct)
}

// AdminProductsPUT is
func AdminProductsPUT(w http.ResponseWriter, r *http.Request) (int, error) {
	return AdminGenericPUT(w, r, new(models.Product), models.UpdateAll)
}
