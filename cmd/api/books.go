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

type JsonValidator struct {
	Errors map[string]string
}

func (jv *JsonValidator) AddError(key, message string) {
	if _, ok := jv.Errors[key]; !ok {
		jv.Errors[key] = message
	}
}
func (jv *JsonValidator) Check(ok bool, key, message string) {
	if !ok {
		jv.AddError(key, message)
	}
}
func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title     string   `json:"title" validate:"required,max=56"`
		Published int32    `json:"published" validate:"required"`
		Pages     int32    `json:"pages" validate:"required,gt=0"`
		Genres    []string `json:"genres" validate:"required"`
	}
	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	jv := JsonValidator{
		Errors: make(map[string]string),
	}
	err = validate.Struct(input)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			fmt.Println(err)
			for _, e := range validateErrs {
				fmt.Println(e.Tag(), e.Param())
				switch {
				case e.Tag() == "required":
					jv.AddError(strings.ToLower(e.Field()), "must be provided")
				case e.Tag() == "max":
					jv.AddError(strings.ToLower(e.Field()), fmt.Sprintf("above the character limit: %v", e.Param()))
				case e.Tag() == "gt":
					jv.AddError(strings.ToLower(e.Field()), "must be provided")
				}

				// jv.Check(e.Tag() == "required", strings.ToLower(e.Field()), "must be provided")
				// jv.Check(e.Tag() == "max", strings.ToLower(e.Field()), fmt.Sprintf("the character limit: %v", e.Param()))
				// jv.Check(e.Tag() == "gt", strings.ToLower(e.Field()), fmt.Sprintf("below the required value: %v", e.Param()))
				// fmt.Println(e)
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
