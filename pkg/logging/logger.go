package logging

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
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
	data := map[string]logrus.Level{
		"DEBUG":   logrus.DebugLevel,
		"INFO":    logrus.InfoLevel,
		"WARNING": logrus.WarnLevel,
		"ERROR":   logrus.ErrorLevel,
		"FATAL":   logrus.FatalLevel,
	}

	Logger = logrus.New()
	Logger.SetReportCaller(true) // this is set to true so that it show correctly where the actual file is

	// set logger level
	var loggerLevel logrus.Level
	if value, ok := data[level]; ok {
		loggerLevel = value
	} else {
		loggerLevel = logrus.DebugLevel
	}
	Logger.SetLevel(loggerLevel)

	Logger.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	Logger.AddHook(&LineHook{})
}
