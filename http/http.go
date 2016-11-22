package http

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
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

			// TODO: show page
			w.Write([]byte(msg.Message))
		}()

		s, err := c.Services.Session.Get(w, r)
		if err != nil {
			return
		}

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

// page is the type that contains the information that goes into the page
type page struct {
	IsLoggedIn  bool
	BaseAddress string
	Data        interface{}
	InviteOnly  bool
	Session     struct {
		FirstName string
		LastName  string
		Email     string
		IsAdmin   bool
		Invites   int
		Credit    int
		Referral  string
	}
}

// Render renders an HTML response and send it to the client based on the
// chosen templates
func Render(w http.ResponseWriter, c *fest.Config, s *fest.Session, data interface{}, templates ...string) (int, error) {
	if strings.HasPrefix(templates[0], "admin/") {
		templates = append(templates, "admin/base")
	} else {
		templates = append(templates, "base")
	}

	var tpl *template.Template

	funcs := template.FuncMap{
		"MD5": func(s string) string {
			hasher := md5.New()
			hasher.Write([]byte(s))
			return hex.EncodeToString(hasher.Sum(nil))
		},
		"DisplayCents": displayCents,
	}

	// For each template, add it to the the tpl variable
	for i := range templates {
		// Get the template from the assets
		page, err := ioutil.ReadFile(filepath.Clean(c.Templates + templates[i] + ".tmpl"))

		// Check if there is some error. If so, the template doesn't exist
		if err != nil {
			c.Logger.Print(err)
			return http.StatusInternalServerError, err
		}

		// If it's the first iteration, creates a new template and add the
		// functions map
		if i == 0 {
			tpl, err = template.New(templates[i]).Funcs(funcs).Parse(string(page))
		} else {
			tpl, err = tpl.Parse(string(page))
		}

		if err != nil {
			c.Logger.Print(err)
			return http.StatusInternalServerError, err
		}
	}

	p := &page{
		IsLoggedIn:  s.Logged,
		Data:        data,
		BaseAddress: c.BaseAddress,
		InviteOnly:  c.InviteOnly,
	}

	// Refresh user information
	if p.IsLoggedIn {
		p.Session.FirstName = s.User.FirstName
		p.Session.LastName = s.User.LastName
		p.Session.Email = s.User.Email
		p.Session.Referral = s.User.Referral
		p.Session.IsAdmin = s.User.Admin
		p.Session.Credit = s.User.Credit
		p.Session.Invites = s.User.Invites
	}

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = buf.WriteTo(w)
	return 0, nil
}

func displayCents(cents int) string {
	price := strconv.Itoa(cents)

	if len(price) == 1 {
		price = "0.0" + price
	} else if len(price) == 2 {
		price = "0." + price
	} else {
		cents := price[len(price)-2:]
		price = price[0:len(price)-2] + "." + cents
	}

	return price
}

// StaticHandler ...
func StaticHandler(templates ...string) FestHandler {
	return func(w http.ResponseWriter, r *http.Request, c *fest.Config) (int, error) {
		s := r.Context().Value("session").(*fest.Session)

		return Render(w, c, s, nil, templates...)
	}
}
