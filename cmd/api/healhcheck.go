package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"version":     version,
		"environment": app.config.env,
	}
	resp, err := json.Marshal(data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	resp = append(resp, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)

}
