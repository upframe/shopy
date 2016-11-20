package http

import (
	"database/sql"
	"net/http"

	"github.com/upframe/fest"
)

type settings struct {
	User    *fest.User
	BaseURL string
}

// SettingsGet ...
func SettingsGet(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
	s := r.Context().Value("session").(*fest.SessionCookie)

	user, err := c.Services.User.Get(s.UserID)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return Render(w, c, s, settings{
		User:    user,
		BaseURL: c.BaseAddress,
	}, "settings")
}
