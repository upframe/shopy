package pages

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"net/http"
	"strconv"
	"strings"

	"github.com/upframe/fest/models"
)

// TODO: REVIEW THIS FILE TO WORK WITH NEW CHANGES

type order struct {
	Cart      models.Cart
	Promocode *models.Promocode
	Credits   int
	Total     float32
}

func init() {
	gob.Register(order{})
	gob.Register(models.Promocode{})
}

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
		// Checks if there are any products in the cart. If there aren't any
		// products, redirect to the cart.
		if len(s.Values["Cart"].(models.Cart).Products) == 0 {
			return Redirect(w, r, "/cart")
		}

		return RenderHTML(w, s, s.Values["Cart"], "checkout-discounts")
	case "pay":
		// Checks if there are any products in the order. If there aren't any
		// products, redirect to the cart.
		if len(s.Values["Order"].(order).Cart.Products) == 0 {
			return Redirect(w, r, "/cart")
		}

		return RenderHTML(w, s, s.Values["Order"], "checkout-pay")
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
		// Creates a new order, adds the Cart of the session, the Promocode,
		// and the credits
		o := order{}
		o.Cart = s.GetCart()
		o.Total = float32(o.Cart.GetTotal())
		o.Credits, err = strconv.Atoi(r.FormValue("credits"))
		if err != nil {
			return http.StatusInternalServerError, err
		}

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

			o.Promocode = generic.(*models.Promocode)
		}

		// Checks if the user has the requested amount of credits
		if s.User.Credit < o.Credits {
			return http.StatusInternalServerError, err
		}

		// Makes the discount
		switch {
		case o.Promocode.Discount >= 1:
			o.Total -= float32(o.Promocode.Discount)
		case o.Promocode.Discount > 0 && o.Promocode.Discount < 1:
			o.Total *= float32(o.Promocode.Discount)
		}

		// Discounts the credits
		o.Total -= float32(o.Credits)

		if o.Total < 0 {
			o.Total = 0
		}

		// Saves the cookie
		s.Values["Order"] = o
		err = s.Save(r, w)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		return http.StatusOK, nil
	case "pay":

		return http.StatusOK, nil
	default:
		return http.StatusNotFound, nil
	}
}

func (o order) creditCardPayment(token string) error {

	/* stripe.Key = "sk_test_GQnowjvTXpLIOMMgceunDKwZ"

	chargeParams := &stripe.ChargeParams{
	  Amount: 2000,
	  Currency: "eur",
	  Desc: "Charge for abigail.thomas@example.com",
	}
	chargeParams.SetSource("tok_192KnBBGLkCZY8NfjR9aCUmE")
	ch, err := charge.New(chargeParams)
	*/

	return nil
}

func (o order) payPalPayment(token string) error {
	return nil
}

// ValidatePromocode validates a promocode and returns the discount amount
// if it exists.
func ValidatePromocode(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	if !s.IsLoggedIn() {
		return http.StatusUnauthorized, nil
	}

	code := new(bytes.Buffer)
	code.ReadFrom(r.Body)

	p, err := models.GetPromocodeByCode(string(code.Bytes()))
	promocode := p.(*models.Promocode)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, nil
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Write([]byte(strconv.Itoa(promocode.Discount)))
	return http.StatusOK, nil
}
