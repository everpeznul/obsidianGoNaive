package http

import (
	"log/slog"
)

var obsiLog *slog.Logger

func HttpSetLog(log *slog.Logger) {

	obsiLog = log
}
