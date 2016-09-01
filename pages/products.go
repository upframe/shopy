package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func AdminProductsGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminProductsPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminProductsDELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminProductsPUT(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}
