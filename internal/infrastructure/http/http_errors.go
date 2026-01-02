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
		// Клиент ушёл или запрос отменён — не пишем ответ

		obsiLog.Debug("WriteError", "error", "context canceled")
		return true

	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "Timeout", http.StatusGatewayTimeout)

		obsiLog.Debug("WriteError", "error", "context deadline")
		return true

	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, "Note wasn't found", http.StatusNotFound)

		obsiLog.Debug("WriteError", "error", "sql no rows")
		return true

	default:
		http.Error(w, "Internal error", http.StatusInternalServerError)

		obsiLog.Debug("WriteError", "error", "not typed")
		return true
	}
}
