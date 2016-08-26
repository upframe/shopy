package login

import "net/http"

func ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusOK, nil
}
