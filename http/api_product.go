package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/upframe/shopy"
	"github.com/gorilla/mux"
)

// APIProductGet ...
func APIProductGet(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
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
func APIProductPost(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	p := &shopy.Product{}

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

	http.Redirect(w, r, "/api/products/"+strconv.Itoa(p.ID), http.StatusSeeOther)
	return 0, nil
}

// APIProductPatch ...
func APIProductPatch(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusNotFound, nil
	}

	p := &shopy.Product{}

	// Get the JSON information
	rawBuffer := new(bytes.Buffer)
	rawBuffer.ReadFrom(r.Body)

	// Parses the JSON into the promocode object and checks for errors
	err = json.Unmarshal(rawBuffer.Bytes(), p)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if p.ID == 0 {
		p.ID = id
	}

	fields, err := topLevelKeys(rawBuffer.Bytes())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = c.Services.Product.Update(p, fields...)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// APIProductDelete  ...
func APIProductDelete(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
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
