package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/tellmeac/expenses/internal/app/app"
	conf "github.com/tellmeac/expenses/internal/app/config"
	"github.com/tellmeac/expenses/internal/app/storage"
	"github.com/tellmeac/expenses/internal/pkg/config"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfgLoader := config.PrepareLoader(config.WithConfigPath("./config.yaml"))

	logger.Info("Parsing config")
	cfg, err := conf.ParseConfig(cfgLoader)
	if err != nil {
		log.Fatalf("Failed to parse config: %s", err)
	}

	logger.Info("Running migrations")
	if err := RunMigrations(ctx, cfg); err != nil {
		log.Fatalf("Migrations failed: %s", err)
	}

	logger.Info("Initializing storage")
	s, err := storage.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Storage init failed: %s", err)
	}

	logger.Info("Initializing app")
	application := app.New(s)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/api/v1/expenses", application.AddExpense)
	r.Get("/api/v1/expenses", application.ListExpenses)

	logger.With("port", cfg.ListenPort).Info("Starting server")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.ListenPort), r))
}

func RunMigrations(ctx context.Context, cfg *conf.Config) error {
	db, err := sql.Open("pgx", cfg.DatabaseConnection)
	if err != nil {
		return err
	}

	return goose.UpContext(ctx, db, "./migrations")
}
