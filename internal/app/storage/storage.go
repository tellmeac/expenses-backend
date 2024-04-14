package storage

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/tellmeac/expenses/internal/app/config"
	"github.com/tellmeac/expenses/internal/app/storage/postgres"
)

type Storage struct {
	db       *sqlx.DB
	Expenses *postgres.Expenses
}

func New(_ context.Context, c *config.Config) (*Storage, error) {
	db, err := sqlx.Connect("pgx", c.DatabaseConnection)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db:       db,
		Expenses: postgres.NewExpenses(db),
	}, nil
}
