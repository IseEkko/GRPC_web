package initialize

import "go.uber.org/zap"

func Logger() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}
