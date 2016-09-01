package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func AdminOrdersGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminOrdersPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminOrdersDELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminOrdersPUT(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}
