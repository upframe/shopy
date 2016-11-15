package http

import (
	"encoding/json"
	"net/http"
)

// TODO: How can we reduce the API code? Delete, Put and Post functions
// are basically the same.

func apiPrint(w http.ResponseWriter, o interface{}) (int, error) {
	data, err := json.MarshalIndent(o, "", "\t")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Write(data)
	return 0, nil
}

func topLevelKeys(data []byte) ([]string, error) {
	// a map container to decode the JSON structure into
	c := make(map[string]interface{})

	// unmarschal JSON
	error := json.Unmarshal(data, &c)

	// panic on error
	if error != nil {
		return []string{}, error
	}

	// a string slice to hold the keys
	k := make([]string, len(c))

	// iteration counter
	i := 0

	// copy c's keys into k
	for s := range c {
		k[i] = s
		i++
	}

	return k, nil
}
