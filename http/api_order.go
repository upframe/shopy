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

// APIOrderGet  ...
func APIOrderGet(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusNotFound, nil
	}

	o, err := c.Services.Order.Get(id)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, nil
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	s := r.Context().Value("session").(*shopy.Session)
	if !s.User.Admin && s.User.ID != o.User.ID {
		return http.StatusForbidden, nil
	}

	return apiPrint(w, o)
}

// APIOrderPost ...
func APIOrderPost(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	o := &shopy.Order{}

	// Get the JSON information
	rawBuffer := new(bytes.Buffer)
	rawBuffer.ReadFrom(r.Body)

	// Parses the JSON into the promocode object and checks for errors
	err := json.Unmarshal(rawBuffer.Bytes(), o)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = c.Services.Order.Create(o)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/api/orders/"+strconv.Itoa(o.ID), http.StatusSeeOther)
	return 0, nil
}

// APIOrderPatch ...
func APIOrderPatch(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusNotFound, nil
	}

	o := &shopy.Order{}

	// Get the JSON information
	rawBuffer := new(bytes.Buffer)
	rawBuffer.ReadFrom(r.Body)

	// Parses the JSON into the promocode object and checks for errors
	err = json.Unmarshal(rawBuffer.Bytes(), o)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if o.ID == 0 {
		o.ID = id
	}

	fields, err := topLevelKeys(rawBuffer.Bytes())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = c.Services.Order.Update(o, fields...)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// APIOrderDelete  ...
func APIOrderDelete(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusNotFound, nil
	}

	err = c.Services.Order.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, nil
}
