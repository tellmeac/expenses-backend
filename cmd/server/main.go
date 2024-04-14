package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/tellmeac/expenses/internal/expense"
	"github.com/tellmeac/expenses/internal/pkg/server"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pressly/goose/v3"
	conf "github.com/tellmeac/expenses/internal/config"
	"github.com/tellmeac/expenses/internal/pkg/config"
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

	db, err := sqlx.ConnectContext(ctx, "pgx", cfg.DatabaseConnection)
	if err != nil {
		log.Fatalf("Connect database: %s", err)
	}

	logger.Info("Running migrations")
	if err = RunMigrations(ctx, db); err != nil {
		log.Fatalf("Migrations failed: %s", err)
	}

	logger.Info("Initializing expenses")
	expenses := expense.New(db)

	srv := server.DefaultServer()
	srv.Addr = fmt.Sprintf(":%s", cfg.ListenPort)

	r := chi.NewRouter()
	srv.Handler = r

	r.Use(middleware.Logger)
	r.Post("/api/v1/expenses", expenses.AddExpense)
	r.Get("/api/v1/expenses", expenses.ListExpenses)
	r.Delete("/api/v1/expenses", expenses.DeleteExpenses)

	logger.With("port", cfg.ListenPort).Info("Starting server")

	log.Fatal(srv.ListenAndServe())
}

func RunMigrations(ctx context.Context, db *sqlx.DB) error {
	return goose.UpContext(ctx, db.DB, "./migrations")
}
