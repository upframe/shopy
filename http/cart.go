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

	cart, err := c.Services.Cart.Get(w, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = cart.FillProducts(c.Services.Product)
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

	cart, err := c.Services.Cart.Get(w, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if cart.Locked {
		return http.StatusForbidden, nil
	}
	if _, ok := cart.RawList[id]; ok {
		// If the Product is already in the cart, increment the quantity
		// Notice that in order for this to work, we have to use pointers
		// (check line 20) and not "normal" values
		cart.RawList[id]++
	} else {
		// Otherwise, we just create a new Cart item, with the product
		cart.RawList[id] = 1
	}

	err = c.Services.Cart.Save(w, cart)
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
	cart, err := c.Services.Cart.Get(w, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if cart.Locked {
		return http.StatusForbidden, nil
	}
	if _, ok := cart.RawList[id]; ok {
		if cart.RawList[id]-1 == 0 {
			delete(cart.RawList, id)
		} else {
			cart.RawList[id]--
		}
	}

	err = c.Services.Cart.Save(w, cart)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
