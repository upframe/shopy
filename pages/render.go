package pages

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// BaseAddress is the base URL of the website
var BaseAddress = "http://upframe.xyz"

// RenderHTML renders an HTML response and send it to the client based on the
// choosen templates
func RenderHTML(w http.ResponseWriter, data interface{}, templates ...string) (int, error) {
	templates = append(templates, "base")
	var tpl *template.Template

	// For each template, add it to the the tpl variable
	for i := range templates {
		// Get the template from the assets
		page, err := ioutil.ReadFile("templates/" + templates[i] + ".tmpl")

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

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, data)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = buf.WriteTo(w)
	return http.StatusOK, nil
}
