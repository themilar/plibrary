package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", app.config.port),
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)
	return srv.ListenAndServe()
}
