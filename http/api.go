package http

import (
	"encoding/json"
	"net/http"
)

func apiPrint(w http.ResponseWriter, o interface{}) (int, error) {
	data, err := json.MarshalIndent(o, "", "\t")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Write(data)
	return 0, nil
}
