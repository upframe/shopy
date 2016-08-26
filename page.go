package upframe

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Info  interface{}
}

func (p *Page) Render(w http.ResponseWriter, templates ...string) (int, error) {
	templates = append(templates, "actions", "base")
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
	err := tpl.Execute(buf, p.Info)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = buf.WriteTo(w)
	return http.StatusOK, nil
}
