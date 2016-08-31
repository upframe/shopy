package pages

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/hacdias/upframe/models"
)

const itemsPerPage = 20

// AdminPromocodesGET redirects to /admin/promocodes/page/1,
func AdminPromocodesGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	// Redirects the user to the first page if he's on /admin/promocodes
	if r.URL.Path == "/admin/promocodes" {
		return Redirect(w, r, "/admin/promocodes/page/1")
	}

	// If the user is not in any page, return a not found status
	if !strings.HasPrefix(r.URL.Path, "/admin/promocodes/page/") {
		return http.StatusNotFound, nil
	}

	// Gets the number of the page
	page, err := strconv.Atoi(strings.Replace(r.URL.Path, "/admin/promocodes/page/", "", 1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Calculates the offset and gets the promocodes
	offset := (page - 1) * itemsPerPage
	promocodes, err := models.GetPromocodes(offset, itemsPerPage)

	// Checks if there are any promocodes. If we're in the first page, show
	// it anyway so we're able to create new promocodes
	if err == sql.ErrNoRows && page != 1 && len(promocodes) == 0 {
		return http.StatusNotFound, err
	}

	// Checks for other errors
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Show the page with the table
	return RenderHTML(w, s, promocodes, "admin/promocodes")
}

// AdminPromocodesPOST creates a new item
func AdminPromocodesPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusOK, nil
}

// AdminPromocodesDELETE deactivates a promocode
func AdminPromocodesDELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	// Removes the "/admin/promocodes/" part from the URL and converts the integer
	// string into a integer variable. Checks for errors
	id, err := strconv.Atoi(strings.Replace(r.URL.Path, "/admin/promocodes/", "", 1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Gets the promocode and checks if it exists
	promocode, err := models.GetPromocode(id)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	// Checks for additional errors
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Deactivates the promocode and checks for errors
	err = promocode.Deactivate()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// AdminPromocodesPUT changes a promocode
func AdminPromocodesPUT(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusOK, nil
}
