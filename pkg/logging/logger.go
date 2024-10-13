package logging

import (
	"runtime"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *logrus.Logger
	Sugar  *zap.SugaredLogger
)

type LineHook struct{}

// all logger levels can execute this hook
func (hook *LineHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// add file and line when logger in fired
func (hook *LineHook) Fire(entry *logrus.Entry) error {
	_, file, line, _ := runtime.Caller(0)
	entry.Data["file"] = file
	entry.Data["line"] = line
	return nil
}

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
