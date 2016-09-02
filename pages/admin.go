package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/upframe/middleware/models"
)

const itemsPerPage = 50

// AdminPromocodesGET handles the GET request for every /admin/promocodes/... URLs
func AdminPromocodesGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
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

// ORDERS
func AdminOrdersGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminOrdersPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminOrdersDELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminOrdersPUT(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

// USERS
func AdminUsersGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminUsersPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminUsersDELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminUsersPUT(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

// PRODUCTS
func AdminProductsGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminProductsPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminProductsDELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminProductsPUT(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}
