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

// APIPromocodeGet ...
func APIPromocodeGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	var (
		p   *fest.Promocode
		err error
	)
	code := mux.Vars(r)["id"]

	if r.URL.Query().Get("code") == "true" {
		p, err = c.Services.Promocode.GetByCode(code)
	} else {
		s := r.Context().Value("session").(*fest.Session)
		if !s.IsAdmin() {
			return http.StatusForbidden, nil
		}
		var id int
		id, err = strconv.Atoi(code)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		p, err = c.Services.Promocode.Get(id)
	}

	if err == sql.ErrNoRows {
		return http.StatusNotFound, nil
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return apiPrint(w, p)
}

// APIPromocodePost ...
func APIPromocodePost(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	p := &fest.Promocode{}

	// Get the JSON information
	rawBuffer := new(bytes.Buffer)
	rawBuffer.ReadFrom(r.Body)

	// Parses the JSON into the promocode object and checks for errors
	err := json.Unmarshal(rawBuffer.Bytes(), p)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = c.Services.Promocode.Create(p)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/api/promocode/"+strconv.Itoa(p.ID), http.StatusSeeOther)
	return 0, nil
}

// APIPromocodeDelete  ...
func APIPromocodeDelete(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusNotFound, nil
	}

	err = c.Services.Promocode.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, nil
}
