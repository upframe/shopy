package pages

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/hacdias/upframe/models"
)

// AdminGenericGET handles the three types of GET requests
func AdminGenericGET(w http.ResponseWriter, r *http.Request, s *sessions.Session, kind string, fn models.GetGenerics) (int, error) {
	// Redirects the user to the first page if he's on /admin/item.
	if r.URL.Path == "/admin/"+kind {
		return Redirect(w, r, "/admin/"+kind+"/page/1")
	}

	// If the user wants to create a new promocode, redirect to /item#new and the
	// javascript will take of the rest.
	if r.URL.Path == "/admin/"+kind+"/new" {
		return Redirect(w, r, "/admin/"+kind+"#new")
	}

	// Checks if the user is in a table page.
	if !strings.HasPrefix(r.URL.Path, "/admin/"+kind+"/page/") {
		// Gets the number of the item and checks for errors
		id, err := strconv.Atoi(strings.Replace(r.URL.Path, "/admin/"+kind+"/", "", 1))
		if err != nil {
			return http.StatusNotFound, err
		}

		// Calculates the number of the page
		page := int(math.Ceil(float64(id) / float64(itemsPerPage)))
		return Redirect(w, r, "/admin/"+kind+"/page/"+strconv.Itoa(page)+"#"+strconv.Itoa(id))
	}

	// Gets the number of the page.
	page, err := strconv.Atoi(strings.Replace(r.URL.Path, "/admin/"+kind+"/page/", "", 1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Calculates the offset and gets the item.
	offset := (page - 1) * itemsPerPage
	items, err := fn(offset, itemsPerPage, "id")

	// Checks if there are any item. If we're in the first page, show
	// it anyway so we're able to create new item.
	if page != 1 && len(items) == 0 {
		return http.StatusNotFound, err
	}

	// Checks for other errors.
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// Show the page with the table.
	return RenderHTML(w, s, items, "admin/"+kind)
}

// AdminGenericPOST creates a new item
func AdminGenericPOST(w http.ResponseWriter, r *http.Request, item models.Generic) (int, error) {
	// Get the JSON information
	rawBuffer := new(bytes.Buffer)
	rawBuffer.ReadFrom(r.Body)

	// Parses the JSON into the promocode object and checks for errors
	err := json.Unmarshal(rawBuffer.Bytes(), item)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Inserts the promocode into the database and checks for errors
	err = item.Insert()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// AdminGenericDELETE deletes an item
func AdminGenericDELETE(w http.ResponseWriter, r *http.Request, kind string, fn models.GetGeneric) (int, error) {
	// Removes the "/admin/kind/" part from the URL and converts the integer
	// string into a integer variable. Checks for errors
	id, err := strconv.Atoi(strings.Replace(r.URL.Path, "/admin/"+kind+"/", "", 1))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Gets the item and checks if it exists
	item, err := fn(id)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	// Checks for additional errors
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Deactivates the item and checks for errors
	err = item.Deactivate()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// AdminGenericPUT updates an item
func AdminGenericPUT(w http.ResponseWriter, r *http.Request, item models.Generic, fields ...string) (int, error) {
	// Get the JSON information
	rawBuffer := new(bytes.Buffer)
	rawBuffer.ReadFrom(r.Body)

	// Parses the JSON into the item object and checks for errors
	err := json.Unmarshal(rawBuffer.Bytes(), item)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Updates the item into the database and checks for errors
	err = item.Update(fields...)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
