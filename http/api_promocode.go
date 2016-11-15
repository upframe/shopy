package http

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/upframe/fest"
)

// APIPromocodeGet ...
func APIPromocodeGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	var p *fest.Promocode
	code := mux.Vars(r)["id"]
	id, err := strconv.Atoi(code)

	if err != nil {
		p, err = c.Services.Promocode.GetByCode(code)
	} else {
		s := r.Context().Value("session").(*fest.Session)
		if !s.IsAdmin() {
			return http.StatusForbidden, nil
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
