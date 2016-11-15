package http

import (
	"encoding/json"
	"net/http"
)

type apiMessage struct {
	Code    int
	Message string
	Error   error `json:"-"`
}

/*
// NotFoundAPI ...
type NotFoundAPI handler

func (h *NotFoundAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := json.MarshalIndent(&apiMessage{
		Code:    http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	}, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write(data)
} */

func apiErrors(w http.ResponseWriter, r *http.Request, err *apiMessage) {
	if err.Error != nil {
		err.Message = err.Error.Error()
	}

	if err.Error == nil && err.Code != 0 {
		err.Message = http.StatusText(err.Code)
	}

	if err.Code != 0 {
		data, e := json.MarshalIndent(err, "", "\t")
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(err.Code)
		w.Write(data)
	}
}
