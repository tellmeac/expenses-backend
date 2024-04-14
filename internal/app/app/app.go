package app

import "github.com/tellmeac/expenses/internal/app/storage"

func New(s *storage.Storage) *App {
	return &App{Storage: s}
}

type App struct {
	Storage *storage.Storage
}
