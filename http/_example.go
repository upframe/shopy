package http

import (
	"net/http"

	"github.com/upframe/fest"
)

// ExampleHandler ...
type ExampleHandler struct {
	SessionService fest.SessionService
	UserService    fest.UserService
}

func (h *ExampleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
	)

	switch r.Method {
	case http.MethodGet:
		code, err = h.GET(w, r)
	case http.MethodPost:
		code, err = h.POST(w, r)
	default:
		code, err = http.StatusNotImplemented, nil
	}

	checkErrors(w, code, err)
}

// GET ...
func (h *ExampleHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := GetSession(w, r, h.UserService)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// POST ...
func (h *ExampleHandler) POST(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := GetSession(w, r, h.UserService)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}
