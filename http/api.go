package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiMessage struct {
	Code    int
	Message string
}

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
}

func apiErrors(w http.ResponseWriter, r *http.Request, code *int, err error) {
	msg := &apiMessage{
		Code: *code,
	}

	fmt.Println(err)

	if err != nil {
		msg.Message = err.Error()
	}

	if err == nil && *code != 0 {
		msg.Message = http.StatusText(*code)
	}

	if *code != 0 {
		data, err := json.MarshalIndent(msg, "", "\t")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(*code)
		w.Write(data)
	}
}
