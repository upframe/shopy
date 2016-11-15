package http

import (
	"net/http"

	"github.com/upframe/fest"
)

// ExampleGET ...
func ExampleGET(cfg fest.Config, srvc fest.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

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
	defer checkErrors(w, r, code, err)

	switch r.Method {
	case http.MethodGet:
		code, err = h.GET(w, r)
	case http.MethodPost:
		code, err = h.POST(w, r)
	default:
		code, err = http.StatusNotImplemented, nil
	}
}

// GET ...
func (h *ExampleHandler) GET(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := GetSession(w, r, c.Services.User)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// POST ...
func (h *ExampleHandler) POST(w http.ResponseWriter, r *http.Request) (int, error) {
	s, err := GetSession(w, r, c.Services.User)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}
