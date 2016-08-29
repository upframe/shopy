package pages

import "net/http"

// IndexGET handles the GET request for /index page
func IndexGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "index")
}
