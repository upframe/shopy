package pages

import "net/http"

func IndexGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "index")
}
