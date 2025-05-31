package main

import (
	"net/http"
)

func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {
	book, headers := createBook(app, w, r)
	if headers != nil {
		err := app.writeJson(w, http.StatusCreated, envelope{"book": book}, headers)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}

func (app *application) bookDetail(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundErrorResponse(w, r)
		return
	}
	book := getBookDetail(app, w, r, id)
	if book != nil {
		if err = app.writeJson(w, http.StatusAccepted, envelope{"book": book}, nil); err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}

}

func (app *application) bookUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundErrorResponse(w, r)
		return
	}
	book := updateBook(app, w, r, id)
	if book != nil {
		err = app.writeJson(w, http.StatusOK, envelope{"book": book}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}

func (app *application) bookDelete(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundErrorResponse(w, r)
		return
	}
	ok := deleteBook(app, w, r, id)
	if ok {
		err = app.writeJson(w, http.StatusNoContent, envelope{}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}

func (app *application) bookList(w http.ResponseWriter, r *http.Request) {
	books, metadata := getBookList(app, w, r)
	if books != nil {
		err := app.writeJson(w, http.StatusOK, envelope{"books": books, "metadata": metadata}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
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
