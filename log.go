package conman

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		exit(FailedToCreateLogsExitCode)
	}
}

// Log returns a pre configured logger
func Log() *zap.Logger {
	return logger
}
