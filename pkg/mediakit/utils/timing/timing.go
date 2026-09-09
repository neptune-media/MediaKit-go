package timing

import (
	"time"

	"github.com/go-logr/logr"
)

func NewFuncDurationPrinter(start time.Time, msg string) func(logger logr.Logger) {
	return func(logger logr.Logger) {
		end := time.Now()
		duration := end.Sub(start)
		logger.Info(msg, "duration", duration)
	}
}
