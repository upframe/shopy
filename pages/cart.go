package pages

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"database/sql"

	"github.com/gorilla/sessions"
	"github.com/upframe/fest/models"
)

type cartItem struct {
	*models.Product
	Quantity int
}

type cart struct {
	Products map[int]*cartItem
	Total    int
}

func init() {
	gob.Register(cartItem{})
	gob.Register(cart{})
}

// CartGET returns the list of items in the cart
func CartGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return Redirect(w, r, "/login")
	}

	return RenderHTML(w, s, s.Values["Cart"], "cart")
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

	cart := s.Values["Cart"].(cart)

	if val, ok := cart.Products[id]; ok {
		// If the Product is already in the cart, increment the quantity
		// Notice that in order for this to work, we have to use pointers
		// (check line 20) and not "normal" values
		val.Quantity++
	} else {
		// Otherwise, we just create a new Cart item, with the product
		cart.Products[id] = &cartItem{
			Product:  product,
			Quantity: 1,
		}
	}

	// Increments the total
	cart.Total += product.Price

	s.Values["Cart"] = cart
	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// CartDELETE removes a product from the cart
func CartDELETE(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return http.StatusUnauthorized, nil
	}

	id, err := strconv.Atoi(strings.Replace(r.URL.Path, "/cart/", "", 1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	fmt.Println(id)

	return http.StatusOK, nil
}
