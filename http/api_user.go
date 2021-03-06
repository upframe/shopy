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

// APIUserGet  ...
func APIUserGet(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusNotFound, nil
	}

	s := r.Context().Value("session").(*shopy.Session)
	if !s.User.Admin && s.User.ID != id {
		return http.StatusForbidden, nil
	}

	p, err := c.Services.User.Get(id)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, nil
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return apiPrint(w, p)
}

// APICurrentUser ...
func APICurrentUser(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	s := r.Context().Value("session").(*shopy.Session)

	return Redirect(w, r, "/api/users/"+strconv.Itoa(s.User.ID))
}

// APIUserPost  ...
func APIUserPost(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	u := &shopy.User{}

	// Get the JSON information
	rawBuffer := new(bytes.Buffer)
	rawBuffer.ReadFrom(r.Body)

	// Parses the JSON into the promocode object and checks for errors
	err := json.Unmarshal(rawBuffer.Bytes(), u)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = c.Services.User.Create(u)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/api/users/"+strconv.Itoa(u.ID), http.StatusSeeOther)
	return 0, nil
}

// APIUserPatch ...
func APIUserPatch(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusNotFound, nil
	}

	u := &shopy.User{}

	// Get the JSON information
	rawBuffer := new(bytes.Buffer)
	rawBuffer.ReadFrom(r.Body)

	// Parses the JSON into the promocode object and checks for errors
	err = json.Unmarshal(rawBuffer.Bytes(), u)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if u.ID == 0 {
		u.ID = id
	}

	fields, err := topLevelKeys(rawBuffer.Bytes())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for i := range fields {
		if fields[i] == "PasswordHash" || fields[i] == "PasswordSalt" {
			fields = append(fields[:i], fields[i+1:]...)
		}
	}

	err = c.Services.User.Update(u, fields...)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// APIUserDelete  ...
func APIUserDelete(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusNotFound, nil
	}

	err = c.Services.User.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, nil
}
