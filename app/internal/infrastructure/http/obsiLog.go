package http

import (
	"log/slog"
)

var obsiLog *slog.Logger

func SetLog(log *slog.Logger) {

	obsiLog = log
}
