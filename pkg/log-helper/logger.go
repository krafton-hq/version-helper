package log_helper

import "go.uber.org/zap"

func Initialize(debug bool, useJsonLogger bool) {
	var cfg zap.Config
	if useJsonLogger {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}
	if debug {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	zap.S().Debug("Initialize Zap Logger with Debug Level")
}
