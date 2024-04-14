package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CreateExpensesRequest struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Catalog     string `json:"catalog"`
	Description string `json:"description"`
	Cost        int64  `json:"cost"`
}

func (a *App) AddExpense(w http.ResponseWriter, r *http.Request) {
	var rb CreateExpensesRequest
	if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
		BadRequest(w, err)
		return
	}

	date, err := time.Parse("2006-01-02", rb.Date)
	if err != nil {
		BadRequest(w, fmt.Errorf("invalid date format: %s", err))
		return
	}

	err = a.Storage.Expenses.Insert(r.Context(), date, rb.Title, rb.Cost, rb.Description, rb.Catalog)
	if err != nil {
		InternalError(w, err)
		return
	}

	NoContent(w)
}

func BadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":    http.StatusBadRequest,
		"message": err.Error(),
	})
}

func InternalError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":    http.StatusInternalServerError,
		"message": err.Error(),
	})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
