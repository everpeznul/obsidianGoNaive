package domain

import (
	"log/slog"
)

var obsiLog *slog.Logger

func DomainSetLog(log *slog.Logger) {

	obsiLog = log
}
