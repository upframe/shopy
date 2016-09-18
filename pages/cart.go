package pages

import (
	"encoding/gob"
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

	// Initialize our data variable and the map of products.
	data := &cart{
		Products: map[int]*cartItem{},
	}

	for _, id := range s.Values["Cart"].([]int) {
		// Gets the product, checks if it exists and checks for errors.
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

		if val, ok := data.Products[id]; ok {
			// If the Product is already in the cart, increment the quantity
			// Notice that in order for this to work, we have to use pointers
			// (check line 20) and not "normal" values
			val.Quantity++
		} else {
			// Otherwise, we just create a new Cart item, with the product
			data.Products[id] = &cartItem{
				Product:  product,
				Quantity: 1,
			}
		}

		// Increments the total
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
