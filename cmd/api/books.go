package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/themilar/plibrary/internal/models"
)

func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title     string   `json:"title"`
		Published int32    `json:"published"`
		Pages     int32    `json:"pages"`
		Genres    []string `json:"genres"`
	}
	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
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
