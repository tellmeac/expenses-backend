-- +goose Up
alter table public.expenses
    add column if not exists is_deleted boolean not null default false,
    add column if not exists deleted_at timestamptz null;

