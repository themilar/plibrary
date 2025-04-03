package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/themilar/plibrary/internal/models"
)

type JsonValidationError struct {
	Errors map[string]string
}

func (jv *JsonValidationError) AddError(key, message string) {
	if _, ok := jv.Errors[key]; !ok {
		jv.Errors[key] = message
	}
}
func validateDate(fl validator.FieldLevel) bool {
	if fl.Field().Int() > int64(time.Now().Year()) {
		return false
	} else if fl.Field().Int() < 1430 {
		return false
	}
	return true
}
func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title     string   `json:"title" validate:"required,max=56"`
		Published int      `json:"published" validate:"required,publication_date"`
		Pages     int      `json:"pages" validate:"required,gt=0"`
		Genres    []string `json:"genres" validate:"required,unique,gt=0,lt=6"`
	}
	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}
	book := &models.Book{
		Title:     input.Title,
		Published: input.Published,
		Pages:     input.Pages,
		Genres:    input.Genres,
	}

	validationErrors := book.Validate()
	if len(validationErrors) != 0 {
		app.failedValidationErrorResponse(w, r, validationErrors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}
func (app *application) bookDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	book := models.Book{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Is it real",
		Pages:     250,
		Version:   1,
		Genres:    []string{"nonfiction", "biography"},
	}
	if err = app.writeJson(w, http.StatusAccepted, envelope{"book": book}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
