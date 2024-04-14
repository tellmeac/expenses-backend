package app

import (
	"github.com/tellmeac/expenses/internal/app/storage"
)

func New(st *storage.Storage) *App {
	return &App{
		Storage: st,
	}
}

type App struct {
	Storage *storage.Storage
}
