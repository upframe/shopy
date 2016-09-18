package pages

import (
	"net/http"
	"strconv"
	"strings"

	"database/sql"

	"github.com/gorilla/sessions"
	"github.com/upframe/fest/models"
)

type cart struct {
	Products []*models.Product
	Total    int
}

// CartGET returns the list of items in the cart
func CartGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return Redirect(w, r, "/login")
	}

	data := &cart{}

	for _, id := range s.Values["Cart"].([]int) {
		generic, err := models.GetProduct(id)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return http.StatusInternalServerError, err
		}

		product := generic.(*models.Product)
		if product.Deactivated {
			continue
		}

		data.Products = append(data.Products, product)
		data.Total += product.Price
	}

	return RenderHTML(w, s, data, "cart")
}

// CartPOST adds a product to the cart
func CartPOST(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return http.StatusUnauthorized, nil
	}

	id, err := strconv.Atoi(strings.Replace(r.URL.Path, "/cart/", "", -1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	product, err := models.GetProduct(id)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if product.(*models.Product).Deactivated {
		return http.StatusNotFound, nil
	}

	s.Values["Cart"] = append(s.Values["Cart"].([]int), id)
	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
