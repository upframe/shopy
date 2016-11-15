package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// PromocodeAPI ...
type PromocodeAPI handler

func (h *PromocodeAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)
	defer apiErrors(w, r, &code, err)

	switch r.Method {
	case http.MethodGet:
		code, err = h.GET(w, r)
	case http.MethodPost:
		code, err = h.POST(w, r)
	case http.MethodPut:
		code, err = h.PUT(w, r)
	case http.MethodDelete:
		code, err = h.DELETE(w, r)
	default:
		code, err = http.StatusNotImplemented, nil
	}
}

// GET ...
func (h *PromocodeAPI) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	p, err := h.Services.Promocode.Get(id)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	data, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Write(data)
	return 0, nil
}

// POST ...
func (h *PromocodeAPI) POST(w http.ResponseWriter, r *http.Request) (int, error) {

	return 0, nil
}

// PUT ...
func (h *PromocodeAPI) PUT(w http.ResponseWriter, r *http.Request) (int, error) {

	return 0, nil
}

// DELETE ...
func (h *PromocodeAPI) DELETE(w http.ResponseWriter, r *http.Request) (int, error) {

	return 0, nil
}
