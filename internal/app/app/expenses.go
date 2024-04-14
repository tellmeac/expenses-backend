package app

import (
	"encoding/json"
	"fmt"
	"github.com/tellmeac/expenses/internal/app/storage/postgres"
	"github.com/tellmeac/expenses/internal/pkg/types"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func (a *App) AddExpense(w http.ResponseWriter, r *http.Request) {
	var p struct {
		Date        types.Date `json:"date"`
		Title       string     `json:"title"`
		Catalog     string     `json:"catalog"`
		Description string     `json:"description"`
		Cost        int64      `json:"cost"`
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		BadRequest(w, err)
		return
	}

	err := a.Storage.Expenses.Insert(r.Context(), p.Date, p.Title, p.Cost, p.Description, p.Catalog)
	if err != nil {
		InternalError(w, err)
		return
	}

	NoContent(w)
}

func (a *App) ListExpenses(w http.ResponseWriter, r *http.Request) {
	limit, err := GetInt64FromQuery(r, "limit")
	if err != nil {
		BadRequest(w, fmt.Errorf("limit: %s", err))
		return
	}

	offset, err := GetInt64FromQuery(r, "offset")
	if err != nil {
		BadRequest(w, fmt.Errorf("offset: %s", err))
		return
	}

	dateFrom, err := GetDateFromQuery(r, "dateFrom")
	if err != nil {
		BadRequest(w, fmt.Errorf("dateFrom: %s", err))
		return
	}

	dateTo, err := GetDateFromQuery(r, "dateTo")
	if err != nil {
		BadRequest(w, fmt.Errorf("dateTo: %s", err))
		return
	}

	expenses, err := a.Storage.Expenses.List(r.Context(), postgres.ListParams{
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
		Offset:   offset,
		Limit:    limit,
	})
	if err != nil {
		InternalError(w, err)
		return
	}

	RespondJSON(w, map[string]any{
		"values": expenses,
	})
}

func (a *App) DeleteExpenses(w http.ResponseWriter, r *http.Request) {
	idsRaw := r.URL.Query().Get("ids")
	idsStr := strings.Split(idsRaw, ",")
	ids := make([]int64, 0, len(idsStr))
	for _, s := range idsStr {
		id, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			BadRequest(w, fmt.Errorf("ids must contain integers separated by ','"))
			return
		}

		ids = append(ids, id)
	}

	err := a.Storage.Expenses.MarkDeleted(r.Context(), ids...)
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

func RespondJSON(w http.ResponseWriter, val any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(val); err != nil {
		slog.With("error", err).Error("Encode json body")
	}
}

func GetInt64FromQuery(r *http.Request, key string) (int64, error) {
	return strconv.ParseInt(r.URL.Query().Get(key), 10, 0)
}

func GetDateFromQuery(r *http.Request, key string) (types.Date, error) {
	return types.ParseDate(r.URL.Query().Get(key))
}
