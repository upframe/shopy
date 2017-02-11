package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/bruhs/shopy"
)

// APIPromocodeGet ...
func APIPromocodeGet(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	var (
		p   *shopy.Promocode
		err error
	)
	code := mux.Vars(r)["id"]

	if r.URL.Query().Get("code") == "true" {
		p, err = c.Services.Promocode.GetByCode(code)
		if time.Now().Unix() > p.Expires.Unix() || p.Used == p.MaxUsage {
			return http.StatusNotFound, nil
		}

	} else {
		s := r.Context().Value("session").(*shopy.Session)
		if !s.User.Admin {
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
func APIPromocodePost(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	p := &shopy.Promocode{}

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

	http.Redirect(w, r, "/api/promocodes/"+strconv.Itoa(p.ID), http.StatusSeeOther)
	return 0, nil
}

// APIPromocodePatch ...
func APIPromocodePatch(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusNotFound, nil
	}

	p := &shopy.Promocode{}

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

	err = c.Services.Promocode.Update(p, fields...)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// APIPromocodeDelete  ...
func APIPromocodeDelete(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
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
