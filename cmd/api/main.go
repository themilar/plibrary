package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var version = ""

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &application{
		config: cfg,
		logger: logger,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck/", app.healthcheckHandler)
	srv := http.Server{
		Addr:         fmt.Sprintf("%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
