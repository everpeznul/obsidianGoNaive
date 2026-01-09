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

		// Клиент ушёл или запрос отменён — не пишем ответ
	case errors.Is(err, context.Canceled):

		obsiLog.Error("WriteError", "type", "context canceled", "error", err)
		return true

	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "Timeout", http.StatusGatewayTimeout)

		obsiLog.Error("WriteError", "type", "context deadline", "error", err)
		return true

	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, "Note wasn't found", http.StatusNotFound)

		obsiLog.Error("WriteError", "type", "sql no rows", "error", err)
		return true

	default:
		http.Error(w, "Internal error", http.StatusInternalServerError)

		obsiLog.Error("WriteError", "type", "not typed", "error", err)
		return true
	}
}
