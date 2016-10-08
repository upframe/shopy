package pages

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/upframe/fest/models"
)

type order struct {
	Cart      cart
	Promocode *models.Promocode
	Credits   int
	Total     float32
}

func init() {
	gob.Register(order{})
	gob.Register(models.Promocode{})
}

// CheckoutGET handles the GET request for /checkout page
func CheckoutGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return Redirect(w, r, "/login")
	}

	if r.URL.Path == "/checkout" {
		return Redirect(w, r, "/checkout/discounts")
	}

	switch strings.Replace(r.URL.Path, "/checkout/", "", -1) {
	case "discounts":
		return RenderHTML(w, s, s.Values["Cart"], "checkout-discounts")
	case "pay":
		return RenderHTML(w, s, s.Values["Order"], "checkout-pay")
	default:
		return http.StatusNotFound, nil
	}
}

// CheckoutPOST handles the POST request for /checkout page
func CheckoutPOST(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return http.StatusUnauthorized, nil
	}

	if r.URL.Path == "/checkout" {
		return Redirect(w, r, "/checkout/discounts")
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
		o.Cart = s.Values["Cart"].(cart)
		o.Total = float32(o.Cart.Total)
		o.Credits, err = strconv.Atoi(r.FormValue("credits"))
		if err != nil {
			return http.StatusInternalServerError, err
		}

		promocode := r.FormValue("promocode")

		if promocode != "" {
			// Gets the promocode and checks for errors
			generic, err := models.GetPromocodeByCode(r.FormValue("promocode"))
			if err == sql.ErrNoRows {
				return http.StatusNotFound, nil
			}

			if err != nil {
				return http.StatusInternalServerError, err
			}

			o.Promocode = generic.(*models.Promocode)
		}

		// Checks if the user has the requested amount of credits
		if s.Values["Credit"].(int) < o.Credits {
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

		return Redirect(w, r, "/checkout/pay")
	case "pay":
		// TODO
		return http.StatusOK, nil
	default:
		return http.StatusNotFound, nil
	}
}

// ValidatePromocode validates a promocode
func ValidatePromocode(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return Redirect(w, r, "/login")
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
