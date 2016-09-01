package pages

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func AdminUsersGET(w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminUsersPOST(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminUsersDELETE(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}

func AdminUsersPUT(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusNotImplemented, nil
}
