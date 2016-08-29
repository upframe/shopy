package pages

import "net/http"

// StoreGET handles the GET request for /store page
func StoreGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "store")
}
