package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Sugar *zap.SugaredLogger
)

// setup logger
func Setup(level string) {
	data := map[string]zapcore.Level{
		"DEBUG":   zap.DebugLevel,
		"INFO":    zap.InfoLevel,
		"WARNING": zap.WarnLevel,
		"ERROR":   zap.ErrorLevel,
		"FATAL":   zap.FatalLevel,
	}

	config := zap.NewDevelopmentConfig()
	var loggerLevel zapcore.Level
	if value, ok := data[level]; ok {
		loggerLevel = value
	} else {
		loggerLevel = zap.DebugLevel
	}
	config.Level = zap.NewAtomicLevelAt(loggerLevel)

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync() // flushes buffer, if any

	Sugar = logger.Sugar()

}
