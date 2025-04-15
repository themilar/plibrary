package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Error(err.Error(), "request_method", r.Method, "request_url", r.URL.String())
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}
	err := app.writeJson(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
func (app *application) notFoundErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)

}
func (app *application) methodNotAllowedErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
func (app *application) badRequestErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
func (app *application) failedValidationErrorResponse(w http.ResponseWriter, r *http.Request, err map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, err)
}
func (app *application) editConflictErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to complete the update due to a conflict, try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}
