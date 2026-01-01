package http

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
)

func writeError(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false

	case errors.Is(err, context.Canceled):
		// Клиент ушёл или запрос отменён — не пишем ответ.
		return true

	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "Timeout", http.StatusGatewayTimeout)
		return true

	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, "Note wasn't found", http.StatusNotFound)
		return true

	default:
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return true
	}
}
