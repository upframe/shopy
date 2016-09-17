package pages

import (
	"net/http"
	"strconv"
	"strings"

	"database/sql"

	"github.com/gorilla/sessions"
	"github.com/upframe/fest/models"
)

// CartGET returns the list of items in the cart
func CartGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return Redirect(w, r, "/login")
	}

	return RenderHTML(w, s, nil, "cart")
}

// CartPOST adds a product to the cart
func CartPOST(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return Redirect(w, r, "/login")
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
