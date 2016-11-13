package pages

import (
	"net/http"

	"github.com/upframe/fest/models"
)

// AdminProductsGET is
func AdminProductsGET(w http.ResponseWriter, r *http.Request) {
 s, _ := GetSession(w, r)
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
