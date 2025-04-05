package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/themilar/plibrary/internal/models"
)

func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title     string   `json:"title" `
		Published int      `json:"published" `
		Pages     int      `json:"pages" `
		Genres    []string `json:"genres" `
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
	err = app.models.Books.Insert(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))
	err = app.writeJson(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) bookDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundErrorResponse(w, r)
		return
	}
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	if err = app.writeJson(w, http.StatusAccepted, envelope{"book": book}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
