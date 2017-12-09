package http

import (
	"database/sql"
	"net/http"

	"github.com/upframe/shopy"
)

// SettingsGet ...
func SettingsGet(w http.ResponseWriter, r *http.Request, c *shopy.Config) (int, error) {
	s := r.Context().Value("session").(*shopy.Session)

	user, err := c.Services.User.Get(s.User.ID)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Render(w, c, s, user, "settings")
}
