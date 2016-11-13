package pages

import (
	"net/http"

	"github.com/upframe/fest/models"
)

// AdminOrdersGET is
func AdminOrdersGET(w http.ResponseWriter, r *http.Request) {
 s, _ := GetSession(w, r)
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
