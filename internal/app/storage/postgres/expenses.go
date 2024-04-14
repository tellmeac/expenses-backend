package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewExpenses(db *sqlx.DB) *Expenses {
	return &Expenses{db: db}
}

type Expenses struct {
	db *sqlx.DB
}

func (e *Expenses) Insert(
	ctx context.Context,
	date time.Time,
	title string,
	cost int64,
	description string,
	catalog string,
) error {
	const query = `
insert into public.expenses(date, title, cost, description, catalog)
values (:date, :title, :cost, :description, :catalog)
`
	args := map[string]any{
		"date":        date,
		"title":       title,
		"cost":        cost,
		"description": description,
		"catalog":     catalog,
	}

	_, err := e.db.NamedExecContext(ctx, query, args)
	return err
}
