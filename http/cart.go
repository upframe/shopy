package http

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/upframe/fest"
)

// CartHandler ...
type CartHandler handler

func (h *CartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)
	defer checkErrors(w, r, code, err)

	switch r.Method {
	case http.MethodGet:
		code, err = h.GET(w, r)
	default:
		code, err = http.StatusNotImplemented, nil
	}
}

// GET ...
func (h *CartHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	cart, err := s.GetCart(h.Services.Product)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, cart, "cart")
}

// CartItemHandler ...
type CartItemHandler handler

func (h *CartItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)

	switch r.Method {
	case http.MethodGet:
		code, err = http.StatusNotFound, nil
	case http.MethodPost:
		code, err = h.POST(w, r)
	case http.MethodDelete:
		code, err = h.DELETE(w, r)
	default:
		code, err = http.StatusNotImplemented, nil
	}

	checkErrors(w, r, code, err)
}

// POST ...
func (h *CartItemHandler) POST(w http.ResponseWriter, r *http.Request) (int, error) {
	s := r.Context().Value("session").(*fest.Session)

	id, err := strconv.Atoi(strings.Replace(r.URL.Path, "/cart/", "", -1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Gets the product, checks if it exists and checks for errors.
	product, err := h.Services.Product.Get(id)
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

// DELETE ...
func (h *CartItemHandler) DELETE(w http.ResponseWriter, r *http.Request) (int, error) {
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
