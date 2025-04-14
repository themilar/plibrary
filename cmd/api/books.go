package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/themilar/plibrary/internal"
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
func (app *application) bookUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
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
	var input struct {
		Title     *string  `json:"title" `
		Published *int     `json:"published" `
		Pages     *int     `json:"pages" `
		Genres    []string `json:"genres" `
	}
	err = app.readJson(w, r, &input)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}
	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Published != nil {
		book.Published = *input.Published
	}
	if input.Pages != nil {
		book.Pages = *input.Pages
	}
	if input.Genres != nil {
		book.Genres = input.Genres
	}
	validationErrors := book.Validate()
	if len(validationErrors) != 0 {
		app.failedValidationErrorResponse(w, r, validationErrors)
		return
	}
	err = app.models.Books.Update(book)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJson(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) bookDelete(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundErrorResponse(w, r)
		return
	}
	err = app.models.Books.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJson(w, http.StatusNoContent, envelope{}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) bookList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		internal.Filters
	}
	qs := r.URL.Query()
	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	filterTypeErrors := map[string]string{}
	p := app.checkEmptyStrings(qs.Get("page"), "1")
	page, err := strconv.Atoi(p)
	if err != nil {
		filterTypeErrors["page"] = "must be an integer"
	}
	input.Filters.Page = page
	s := app.checkEmptyStrings(qs.Get("size"), "12")
	size, err := strconv.Atoi(s)
	if err != nil {
		filterTypeErrors["size"] = "must be an integer"
	}
	input.Filters.Size = size
	input.Filters.Sort = app.readString(qs, "sort", "id")
	filterErrors := internal.ValidateFilters(input.Filters, filterTypeErrors)
	if len(filterErrors) > 0 {
		app.failedValidationErrorResponse(w, r, filterErrors)
		return
	}
	books, metadata, err := app.models.Books.All(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusOK, envelope{"books": books, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) bookSearch(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	title := app.readString(qs, "q", "")
	books, err := app.models.Books.FullTextSearch(title)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusOK, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
