package pages

import "net/http"

// CartGET handles the GET request for /cart page
func CartGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "cart")
}
