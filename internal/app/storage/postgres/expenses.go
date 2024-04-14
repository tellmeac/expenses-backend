package postgres

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/tellmeac/expenses/internal/pkg/types"
)

type Expense struct {
	ID          int64  `db:"id"`
	Cost        int64  `db:"cost"`
	Date        string `db:"date"`
	Title       string `db:"title"`
	Catalog     string `db:"catalog"`
	Description string `db:"description"`
}

func NewExpenses(db *sqlx.DB) *Expenses {
	return &Expenses{db: db}
}

type Expenses struct {
	db *sqlx.DB
}

func (e *Expenses) Insert(
	ctx context.Context,
	date types.Date,
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
		"date":        date.Time(),
		"title":       title,
		"cost":        cost,
		"description": description,
		"catalog":     catalog,
	}

	_, err := e.db.NamedExecContext(ctx, query, args)
	return err
}

type ListParams struct {
	DateFrom, DateTo *types.Date
	Offset, Limit    int64
}

func (e *Expenses) List(ctx context.Context, p ListParams) ([]Expense, error) {
	if p.Limit == 0 {
		return []Expense{}, nil
	}

	qb := squirrel.Select("id", "cost", "date", "title", "catalog", "description").
		PlaceholderFormat(squirrel.Dollar).
		From("public.expenses").
		Limit(uint64(p.Limit)).Offset(uint64(p.Offset))

	if p.DateFrom != nil {
		qb = qb.Where("date >= ?", p.DateFrom.Time())
	}

	if p.DateTo != nil {
		qb = qb.Where("date <= ?", p.DateTo.Time())
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("unexpected sql: %w", err)
	}

	var result []Expense

	return result, e.db.SelectContext(ctx, &result, sql, args...)
}
