package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	var address string
	if app.config.env == "development" {
		address = fmt.Sprintf("localhost:%d", app.config.port)
	} else if app.config.env == "production" {
		address = fmt.Sprintf(":%v", app.config.port)
	}
	srv := &http.Server{
		Addr:         address,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		app.logger.Info("shutting down server", "signal", s.String())
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdownError <- srv.Shutdown(ctx)

	}()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}
	app.logger.Info("stopped server", "addr", srv.Addr)
	return nil
}
