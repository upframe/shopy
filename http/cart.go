package http

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/upframe/fest"
)

// CartHandler ...
type CartHandler struct {
	UserService    fest.UserService
	ProductService fest.ProductService
}

func (h *CartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)

	switch r.Method {
	case http.MethodGet:
		code, err = h.GET(w, r)
	case http.MethodPost:
		code, err = h.POST(w, r)
	case http.MethodDelete:
		code, err = h.DELETE(w, r)
	default:
		code, err = http.StatusNotImplemented, nil
	}

	checkErrors(w, code, err)
}

// GET ...
func (h *CartHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := GetSession(w, r, h.UserService)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !s.IsLoggedIn() {
		return Redirect(w, r, "/login")
	}

	cart, err := s.GetCart(h.ProductService)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return RenderHTML(w, s, cart, "cart")
}

// POST ...
func (h *CartHandler) POST(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := GetSession(w, r, h.UserService)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !s.IsLoggedIn() {
		return http.StatusUnauthorized, fest.ErrNotLoggedIn
	}

	id, err := strconv.Atoi(strings.Replace(r.URL.Path, "/cart/", "", -1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Gets the product, checks if it exists and checks for errors.
	product, err := h.ProductService.Get(id)
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
func (h *CartHandler) DELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := GetSession(w, r, h.UserService)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !s.IsLoggedIn() {
		return http.StatusUnauthorized, nil
	}

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
