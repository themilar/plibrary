package main

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	if app.config.limiter.enabled {
		router.Use(httprate.LimitByIP(app.config.limiter.rpm, time.Minute))
	}
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:9000"},
	}))
	router.NotFound(app.notFoundErrorResponse)
	router.MethodNotAllowed(app.methodNotAllowedErrorResponse)

	router.Get("/v1/healthcheck", app.healthcheckHandler)
	router.Get("/v1/books", app.bookList)
	router.Get("/v1/books/search", app.bookSearch)
	router.Post("/v1/books", app.bookCreate)
	router.Get("/v1/books/{id}", app.bookDetail)
	router.Patch("/v1/books/{id}", app.bookUpdate)
	router.Delete("/v1/books/{id}", app.bookDelete)
	return router
}
