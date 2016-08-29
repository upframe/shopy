package pages

import "net/http"

// CheckoutGET handles the GET request for /checkout page
func CheckoutGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return RenderHTML(w, nil, "checkout")
}
