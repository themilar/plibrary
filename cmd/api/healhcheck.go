package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := envelope{
		"status": "available",
		"system_info": map[string]string{

			"version":     version,
			"environment": app.config.env,
		}}
	resp, err := json.Marshal(data)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
		return
	}
	resp = append(resp, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)

}
