package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/themilar/plibrary/internal"
	"github.com/themilar/plibrary/internal/models"
)

var createInput struct {
	Title     string   `json:"title" `
	Published int      `json:"published" `
	Pages     int      `json:"pages" `
	Genres    []string `json:"genres" `
}

// type updateInput struct {
// 	Title     *string  `json:"title" `
// 	Published *int     `json:"published" `
// 	Pages     *int     `json:"pages" `
// 	Genres    []string `json:"genres" `
// }

var listInput struct {
	Title  string
	Genres []string
	internal.Filters
}

func createBook(app *application, w http.ResponseWriter, r *http.Request) (*models.Book, http.Header) {
	err := app.readJson(w, r, &createInput)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return &models.Book{}, nil
	}
	book := &models.Book{
		Title:     createInput.Title,
		Published: createInput.Published,
		Pages:     createInput.Pages,
		Genres:    createInput.Genres,
	}

	validationErrors := book.Validate()
	if len(validationErrors) != 0 {
		app.failedValidationErrorResponse(w, r, validationErrors)
		return &models.Book{}, nil
	}
	err = app.models.Books.Insert(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return &models.Book{}, nil
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))
	return book, headers
}
func updateBook(app *application, w http.ResponseWriter, r *http.Request, id int64) *models.Book {
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return nil
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
		return nil
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
		return nil
	}
	err = app.models.Books.Update(book)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return nil
	}
	return book
}
func getBookDetail(app *application, w http.ResponseWriter, r *http.Request, id int64) *models.Book {
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return nil
	}
	return book
}
func deleteBook(app *application, w http.ResponseWriter, r *http.Request, id int64) bool {
	err := app.models.Books.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return false
	}
	return true

}
func getBookList(app *application, w http.ResponseWriter, r *http.Request) ([]*models.Book, *internal.PaginationMetadata) {
	qs := r.URL.Query()
	listInput.Title = app.readString(qs, "title", "")
	listInput.Genres = app.readCSV(qs, "genres", []string{})
	filterTypeErrors := map[string]string{}
	p := app.checkEmptyStrings(qs.Get("page"), "1")
	page, err := strconv.Atoi(p)
	if err != nil {
		filterTypeErrors["page"] = "must be an integer"
	}
	listInput.Filters.Page = page
	s := app.checkEmptyStrings(qs.Get("size"), "12")
	size, err := strconv.Atoi(s)
	if err != nil {
		filterTypeErrors["size"] = "must be an integer"
	}
	listInput.Filters.Size = size
	listInput.Filters.Sort = app.readString(qs, "sort", "id")
	filterErrors := internal.ValidateFilters(listInput.Filters, filterTypeErrors)
	if len(filterErrors) > 0 {
		app.failedValidationErrorResponse(w, r, filterErrors)
		return nil, nil
	}
	books, metadata, err := app.models.Books.All(listInput.Title, listInput.Genres, listInput.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return nil, nil
	}
	return books, metadata
}
