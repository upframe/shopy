package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/hacdias/upframe/models"
)

// StoreGET handles the GET request for /store page
func StoreGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	models.DeleteByID(1)

	return RenderHTML(w, nil, "store")
}
