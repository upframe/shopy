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

	cookie, err := ReadCartCookie(w, r, c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	cart, err := cookie.GetCart(c.Services.Product)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Render(w, c, s, cart, "cart")
}

// CartItemPost ...
func CartItemPost(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
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

	cart, err := ReadCartCookie(w, r, c)
	if err != nil {
		return http.StatusInternalServerError, err
	}
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

	err = SetCartCookie(w, c, cart)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// CartItemDelete ...
func CartItemDelete(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	id, err := strconv.Atoi(strings.Replace(r.URL.Path, "/cart/", "", 1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Remove one item of this type from the cart
	cart, err := ReadCartCookie(w, r, c)
	if err != nil {
		return http.StatusInternalServerError, err
	}
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

	err = SetCartCookie(w, c, cart)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
