package pages

import (
	"net/http"

	"github.com/upframe/fest/models"
)

// IndexGET handles the GET request for /index page
func IndexGET(w http.ResponseWriter, r *http.Request, s *models.Session) (int, error) {
	return RenderHTML(w, s, nil, "index")
}
