package upframe

import "net/http"

type Page struct {
}

func (p *Page) Render(w http.ResponseWriter, templates ...string) (int, error) {

	return http.StatusOK, nil
}
