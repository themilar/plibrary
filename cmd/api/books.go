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
	fmt.Fprintln(w, "create a new book")
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
	if err = app.writeJson(w, http.StatusAccepted, book); err != nil {
		app.logger.Print(err)
		http.Error(w, "the server encountered a problem", http.StatusInternalServerError)
	}

}
