package http

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/bruhs/shopy"
)

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
func Render(w http.ResponseWriter, c *shopy.Config, s *shopy.Session, data interface{}, templates ...string) (int, error) {
	if strings.HasPrefix(templates[0], "admin/") {
		templates = append(templates, "admin/base")
	} else {
		templates = append(templates, "base")
	}

	funcs := template.FuncMap{
		"MD5": func(s string) string {
			hasher := md5.New()
			hasher.Write([]byte(s))
			return hex.EncodeToString(hasher.Sum(nil))
		},
		"DisplayCents": shopy.DisplayCents,
	}

	for i := range templates {
		templates[i] = filepath.Clean(c.Templates + templates[i] + ".tmpl")
	}

	tpl, err := template.New("main").Funcs(funcs).ParseFiles(templates...)
	if err != nil {
		c.Logger.Print(err)
		return http.StatusInternalServerError, err
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
	err = tpl.ExecuteTemplate(buf, "base.tmpl", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = buf.WriteTo(w)
	return 0, err
}
