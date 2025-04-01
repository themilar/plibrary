package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
		Published int32    `json:"published" validate:"required,publication_date"`
		Pages     int32    `json:"pages" validate:"required,gt=0"`
		Genres    []string `json:"genres" validate:"required,unique,gt=0,lt=6"`
	}
	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("publication_date", validateDate)
	jv := JsonValidationError{
		Errors: make(map[string]string),
	}
	err = validate.Struct(input)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				fmt.Println(e.Tag(), e.Param())
				switch {
				case e.Tag() == "required":
					jv.AddError(strings.ToLower(e.Field()), "must be provided")
				case e.Tag() == "max":
					jv.AddError(strings.ToLower(e.Field()), fmt.Sprintf("above the character limit: %v", e.Param()))
				case e.Tag() == "gt":
					jv.AddError(strings.ToLower(e.Field()), "must be provided")
				case e.Tag() == "lt":
					jv.AddError(strings.ToLower(e.Field()), "must not exceed 5 items")
				case e.Tag() == "publication_date":
					jv.AddError(strings.ToLower(e.Field()), fmt.Sprintf("publication date cannot exceed the range: 1430-%v", time.Now().Year()))
				case e.Tag() == "unique":
					jv.AddError(strings.ToLower(e.Field()), "cannot contain duplicate genres")
				}
			}
			app.failedValidationErrorResponse(w, r, jv.Errors)
		}
		// from here you can create your own error messages in whatever language you wish
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
