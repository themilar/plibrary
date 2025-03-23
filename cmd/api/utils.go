package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) writeJson(w http.ResponseWriter, status int, data interface{}) error {
	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp = append(resp, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(resp)
	return nil
}
