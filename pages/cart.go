package pages

import "net/http"

func CartGET(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusOK, nil
}
