package pages

import (
	"bytes"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/upframe/fest/models"
)

// CheckoutGET handles the GET request for /checkout page
func CheckoutGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if !IsLoggedIn(s) {
		return Redirect(w, r, "/login")
	}

	return RenderHTML(w, s, s.Values["Cart"], "checkout")
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
