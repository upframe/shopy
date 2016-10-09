package pages

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/upframe/fest/models"
)

// CartGET returns the list of items in the cart
func CartGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	if !s.IsLoggedIn() {
		return Redirect(w, r, "/login")
	}

	cart, err := s.GetCart()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, cart, "cart")
}

// CartPOST adds a product to the cart
func CartPOST(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	if !s.IsLoggedIn() {
		return http.StatusUnauthorized, nil
	}

	id, err := strconv.Atoi(strings.Replace(r.URL.Path, "/cart/", "", -1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Gets the product, checks if it exists and checks for errors.
	generic, err := models.GetProduct(id)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	product := generic.(*models.Product)
	if product.Deactivated {
		return http.StatusNotFound, err
	}

	cart := s.Values["Cart"].(map[int]int)

	if _, ok := cart[id]; ok {
		// If the Product is already in the cart, increment the quantity
		// Notice that in order for this to work, we have to use pointers
		// (check line 20) and not "normal" values
		cart[id]++
	} else {
		// Otherwise, we just create a new Cart item, with the product
		cart[id] = 1
	}

	s.Values["Cart"] = cart
	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// CartDELETE removes a product from the cart
func CartDELETE(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	if !s.IsLoggedIn() {
		return http.StatusUnauthorized, nil
	}

	id, err := strconv.Atoi(strings.Replace(r.URL.Path, "/cart/", "", 1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Remove one item of this type from the cart
	cart := s.Values["Cart"].(map[int]int)
	if _, ok := cart[id]; ok {
		if cart[id]-1 == 0 {
			delete(cart, id)
		} else {
			cart[id]--
		}
	}

	s.Values["Cart"] = cart

	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
