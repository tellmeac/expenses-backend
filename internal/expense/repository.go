package expense

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/tellmeac/expenses/internal/pkg/types"
)

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *sqlx.DB
}

type Expense struct {
	ID          int64      `db:"id" json:"id"`
	Cost        int64      `db:"cost" json:"cost"`
	Date        types.Date `db:"date" json:"date"`
	Title       string     `db:"title" json:"title"`
	Catalog     string     `db:"catalog" json:"catalog"`
	Description string     `db:"description" json:"description"`
	IsDeleted   bool       `db:"is_deleted" json:"is_deleted"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}

func (e *Repository) Insert(
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
	IsDeleted        bool
}

func (e *Repository) List(ctx context.Context, p ListParams) ([]Expense, error) {
	if p.Limit == 0 {
		return []Expense{}, nil
	}

	qb := squirrel.Select("id", "cost", "date", "title", "catalog", "description", "is_deleted", "deleted_at").
		PlaceholderFormat(squirrel.Dollar).
		From("public.expenses").
		Where("is_deleted = ?", p.IsDeleted).
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

func (e *Repository) MarkDeleted(ctx context.Context, ids ...int64) error {
	if len(ids) == 0 {
		return nil
	}

	sql, args, err := sqlx.In(`update public.expenses set is_deleted = true, deleted_at = now() where id in (?)`, ids)
	if err != nil {
		return fmt.Errorf("build sql: %w", err)
	}
	sql = e.db.Rebind(sql)

	_, err = e.db.ExecContext(ctx, sql, args...)
	return err
}
