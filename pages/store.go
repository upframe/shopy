package pages

import "net/http"

func StoreGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "store")
}
