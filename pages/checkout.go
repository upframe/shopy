package pages

import "net/http"

func CheckoutGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "checkout")
}
