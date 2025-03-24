package main

import (
	"encoding/json"
	"net/http"
)

// func (app *application) readIDParam(r *http.Request) (int64, error) {
// 	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
// 	if err != nil || id < 1 {
// 		return 0, errors.New("invalid id parameter")
// 	}
// 	return id, nil

// }
func (app *application) writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	resp, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	resp = append(resp, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(resp)
	return nil
}

type envelope map[string]interface{}
