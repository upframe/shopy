package http

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/upframe/fest"
)

type message struct {
	ID      string
	Code    int
	Message string
	Error   error `json:"-"`
}

// FestHandler ...
type FestHandler func(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error)

// Inject ...
func Inject(h FestHandler, c *fest.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			code int
			err  error
		)

		s, err := c.Services.Session.Get(w, r)
		if err != nil {
			return
		}

		defer func() {
			if code == 0 && err == nil {
				return
			}

			msg := &message{Code: code}

			if err != nil {
				msg.Message = err.Error()
			} else {
				msg.Message = http.StatusText(code)
			}

			if code >= 400 {
				t := time.Now()
				msg.ID = t.Format("20060102150405")
			}

			if code >= 400 && err != nil {
				c.Logger.Print(err)
			}

			if code != 0 {
				w.WriteHeader(code)
			}

			if strings.HasPrefix(r.URL.Path, "/api") || r.Method != http.MethodGet {
				data, e := json.MarshalIndent(msg, "", "\t")
				if e != nil {
					c.Logger.Print(e)
					return
				}

				w.Write(data)
				return
			}

			if code == http.StatusNotFound {
				_, err = Render(w, c, s, msg, "404")
				if err != nil {
					c.Logger.Print(err)
					return
				}
				return
			}

			var tpl *template.Template
			tpl, err = template.New("errors").Parse(errorTemplate)
			if err != nil {
				c.Logger.Print(err)
				return
			}

			err = tpl.Execute(w, msg)
			if err != nil {
				c.Logger.Print(err)
			}
		}()

		r = r.WithContext(context.WithValue(r.Context(), "session", s))
		code, err = h(w, r, c)
	}
}

// MustLogin ...
func MustLogin(h FestHandler) FestHandler {
	return func(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
		s := r.Context().Value("session").(*fest.Session)

		if s.Logged {
			return h(w, r, c)
		}

		if r.Method == http.MethodGet && !strings.HasPrefix(r.URL.Path, "/api") {
			return Redirect(w, r, "/login?redirect="+url.QueryEscape(r.URL.Path))
		}

		return http.StatusUnauthorized, fest.ErrNotLoggedIn
	}
}

// MustAdmin ...
func MustAdmin(h FestHandler) FestHandler {
	return MustLogin(func(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
		s := r.Context().Value("session").(*fest.Session)

		if s.User.Admin {
			return h(w, r, c)
		}

		return http.StatusForbidden, nil
	})
}

// Redirect redirects the user to a page
func Redirect(w http.ResponseWriter, r *http.Request, path string) (int, error) {
	http.Redirect(w, r, path, http.StatusTemporaryRedirect)
	return 0, nil
}

// StaticHandler ...
func StaticHandler(templates ...string) FestHandler {
	return func(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
		s := r.Context().Value("session").(*fest.Session)

		return Render(w, c, s, nil, templates...)
	}
}
