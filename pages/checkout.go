package pages

import (
	"bytes"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/upframe/fest/models"
)

// CheckoutGET handles the GET request for /checkout page
func CheckoutGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	if !s.IsLoggedIn() {
		return Redirect(w, r, "/login")
	}

	if r.URL.Path == "/checkout" {
		return Redirect(w, r, "/checkout/discounts")
	}

	switch strings.Replace(r.URL.Path, "/checkout/", "", -1) {
	case "discounts":
		cart, err := s.GetCart()
		if err != nil {
			return http.StatusInternalServerError, err
		}

		return RenderHTML(w, s, cart, "checkout-discounts")
	case "pay":
		data := map[string]interface{}{}

		var err error
		data["Cart"], err = s.GetCart()
		if err != nil {
			return http.StatusInternalServerError, err
		}

		data["Order"] = s.Values["Order"].(models.OrderCookie)
		return RenderHTML(w, s, data, "checkout-pay")
	default:
		return http.StatusNotFound, nil
	}
}

// CheckoutPOST handles the POST request for /checkout page
func CheckoutPOST(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	if !s.IsLoggedIn() {
		return http.StatusUnauthorized, nil
	}

	// Parses the form and checks for errors
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	switch strings.Replace(r.URL.Path, "/checkout/", "", -1) {
	case "discounts":
		return checkoutPOSTDiscount(w, r, s)
	case "pay":
		return checkoutPOSTPay(w, r, s)
	default:
		return http.StatusNotFound, nil
	}
}

func checkoutPOSTDiscount(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	cart, err := s.GetCart()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Lock the cart
	cartCookie := s.Values["Cart"].(models.CartCookie)
	cartCookie.Locked = true

	order := models.OrderCookie{}
	order.Total = cart.GetTotal()

	// Obtain the credits and discount them
	credits := r.FormValue("credits")
	if len(credits) == 0 {
		credits = "0"
	}

	order.Credits, err = strconv.Atoi(credits)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if s.User.Credit < order.Credits || order.Credits > order.Total {
		return http.StatusBadRequest, nil
	}

	order.Total -= order.Credits

	// Gets the promocode
	promocode := r.FormValue("promocode")

	if promocode != "" {
		// Gets the promocode and checks for errors
		var generic models.Generic
		generic, err = models.GetPromocodeByCode(r.FormValue("promocode"))
		if err == sql.ErrNoRows {
			return http.StatusNotFound, nil
		}

		if err != nil {
			return http.StatusInternalServerError, err
		}

		promo := generic.(*models.Promocode)
		order.Promocode.Code = promo.Code

		if promo.Percentage {
			order.Promocode.DiscountAmount = (promo.Discount * order.Total) / 100
			order.Promocode.DiscountAmount = order.Total - order.Promocode.DiscountAmount
		} else {
			order.Promocode.DiscountAmount = promo.Discount
		}

		order.Total -= order.Promocode.DiscountAmount
	}

	// Saves the cookie
	s.Values["Cart"] = cartCookie
	s.Values["Order"] = order
	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func checkoutPOSTPay(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	/*

		// IF: STRIPEs
		params := &stripe.ChargeParams{
			Amount:   1000,
			Currency: currency.USD,
			Card:     &stripe.CardParams{Token: "tok_14dlcYGBoqcjK6A1Th7tPXfJ"},
			Desc:     "Gopher t-shirt",
		}

		ch, err := charge.New(params) */

	return http.StatusNotImplemented, nil
}

// ValidatePromocode validates a promocode and returns the discount amount
// if it exists.
func ValidatePromocode(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	if !s.IsLoggedIn() {
		return http.StatusUnauthorized, errNotLoggedIn
	}

	code := new(bytes.Buffer)
	code.ReadFrom(r.Body)

	p, err := models.GetPromocodeByCode(string(code.Bytes()))
	if err == sql.ErrNoRows {
		return http.StatusNotFound, nil
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Write([]byte(strconv.Itoa(p.(*models.Promocode).Discount)))
	return http.StatusOK, nil
}
