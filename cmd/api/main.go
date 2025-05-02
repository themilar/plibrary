package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/themilar/plibrary/internal/models"
)

var version string = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	limiter struct {
		enabled bool
		rpm     int
	}
}

type application struct {
	config config
	logger *slog.Logger
	models models.Models
}

func main() {
	var cfg config
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading env file")
	}

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")
	flag.BoolVar(&cfg.limiter.enabled, "limitenabled", true, "Enable rate limiter")
	flag.IntVar(&cfg.limiter.rpm, "limitrpm", 50, "rate limiter maximum requests per minute")
	flag.Parse()

	db, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		logger.Error(err.Error())
	}
	defer db.Close()
	logger.Info("Database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
	}

}
