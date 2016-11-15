package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/upframe/fest"
)

// APIProductGet ...
func APIProductGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusNotFound, nil
	}

	p, err := c.Services.Product.Get(id)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, nil
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return apiPrint(w, p)
}

// APIProductPost ...
func APIProductPost(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	p := &fest.Product{}

	// Get the JSON information
	rawBuffer := new(bytes.Buffer)
	rawBuffer.ReadFrom(r.Body)

	// Parses the JSON into the promocode object and checks for errors
	err := json.Unmarshal(rawBuffer.Bytes(), p)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = c.Services.Product.Create(p)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/api/product/"+strconv.Itoa(p.ID), http.StatusSeeOther)
	return 0, nil
}

// APIProductDelete  ...
func APIProductDelete(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusNotFound, nil
	}

	err = c.Services.Product.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, nil
}
