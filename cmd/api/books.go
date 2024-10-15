package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new book")
}
func (app *application) bookDetail(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	fmt.Fprintf(w, "details of the book %d\n", id)
}
