package http

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/upframe/fest"
)

// APIOrderGet  ...
func APIOrderGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return http.StatusNotFound, nil
	}

	p, err := c.Services.Order.Get(id)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, nil
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return apiPrint(w, p)
}
