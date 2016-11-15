package http

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/upframe/fest"
)

// CartGet ...
func CartGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	cart, err := s.GetCart(c.Services.Product)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, c, s, cart, "cart")
}

// CartItemPost ...
func CartItemPost(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	id, err := strconv.Atoi(strings.Replace(r.URL.Path, "/cart/", "", -1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Gets the product, checks if it exists and checks for errors.
	product, err := c.Services.Product.Get(id)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if product.Deactivated {
		return http.StatusNotFound, err
	}

	cart := s.Values["Cart"].(fest.CartCookie)
	if cart.Locked {
		return http.StatusForbidden, nil
	}
	if _, ok := cart.Products[id]; ok {
		// If the Product is already in the cart, increment the quantity
		// Notice that in order for this to work, we have to use pointers
		// (check line 20) and not "normal" values
		cart.Products[id]++
	} else {
		// Otherwise, we just create a new Cart item, with the product
		cart.Products[id] = 1
	}

	s.Values["Cart"] = cart
	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// CartItemDelete ...
func CartItemDelete(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	id, err := strconv.Atoi(strings.Replace(r.URL.Path, "/cart/", "", 1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Remove one item of this type from the cart
	cart := s.Values["Cart"].(fest.CartCookie)
	if cart.Locked {
		return http.StatusForbidden, nil
	}
	if _, ok := cart.Products[id]; ok {
		if cart.Products[id]-1 == 0 {
			delete(cart.Products, id)
		} else {
			cart.Products[id]--
		}
	}

	s.Values["Cart"] = cart

	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
