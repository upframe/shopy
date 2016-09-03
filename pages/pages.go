package pages

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
)

var (
	// BaseAddress is the base URL of the website
	BaseAddress = "http://upframe.xyz"
	// TemplatesPath is the root path of the website
	TemplatesPath string
)

// page is the type that contains the information that goes into the page
type page struct {
	IsLoggedIn bool
	Data       interface{}
	Session    struct {
		FirstName string
		LastName  string
		IsAdmin   bool
	}
}

// RenderHTML renders an HTML response and send it to the client based on the
// choosen templates
func RenderHTML(w http.ResponseWriter, s *sessions.Session, data interface{}, templates ...string) (int, error) {
	templates = append(templates, "base")
	var tpl *template.Template

	// For each template, add it to the the tpl variable
	for i := range templates {
		// Get the template from the assets
		page, err := ioutil.ReadFile(filepath.Clean(TemplatesPath + templates[i] + ".tmpl"))

		// Check if there is some error. If so, the template doesn't exist
		if err != nil {
			log.Print(err)
			return http.StatusInternalServerError, err
		}

		// If it's the first iteration, creates a new template and add the
		// functions map
		if i == 0 {
			tpl, err = template.New(templates[i]).Parse(string(page))
		} else {
			tpl, err = tpl.Parse(string(page))
		}

		if err != nil {
			log.Print(err)
			return http.StatusInternalServerError, err
		}
	}

	p := &page{
		IsLoggedIn: IsLoggedIn(s),
		Data:       data,
	}

	if p.IsLoggedIn {
		p.Session.FirstName = s.Values["FirstName"].(string)
		p.Session.LastName = s.Values["LastName"].(string)
		p.Session.IsAdmin = s.Values["IsAdmin"].(bool)
	}

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = buf.WriteTo(w)
	return http.StatusOK, nil
}

// IsLoggedIn checks if an user is logged in
func IsLoggedIn(s *sessions.Session) bool {
	switch s.Values["IsLoggedIn"].(type) {
	case bool:
		return s.Values["IsLoggedIn"].(bool)
	}

	return false
}

// IsAdmin checks if an user is admin
func IsAdmin(s *sessions.Session) bool {
	switch s.Values["IsAdmin"].(type) {
	case bool:
		return s.Values["IsAdmin"].(bool)
	}

	return false
}

// Redirect redirects the user to a page
func Redirect(w http.ResponseWriter, r *http.Request, path string) (int, error) {
	http.Redirect(w, r, path, http.StatusTemporaryRedirect)
	return http.StatusOK, nil
}
