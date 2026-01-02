package use_case

import (
	"log/slog"
)

var obsiLog *slog.Logger

func UseCaseSetLog(log *slog.Logger) {

	obsiLog = log
}
