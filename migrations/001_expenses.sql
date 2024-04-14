-- +goose Up
create table if not exists public.expenses(
    id serial primary key,
    date date not null,
    title text not null,
    catalog varchar(128) not null,
    cost integer not null,
    description text null
);
