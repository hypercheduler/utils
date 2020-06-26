package log

import (
	"go.uber.org/zap"
	"os"
	"sync"
)

var lock sync.Mutex
var loggerMap map[string]*zap.Logger

var LevelInfo = zap.NewAtomicLevelAt(zap.InfoLevel)
var LevelError = zap.NewAtomicLevelAt(zap.ErrorLevel)
var LevelFatal = zap.NewAtomicLevelAt(zap.FatalLevel)

var defaultConfig *LoggerConfig

type LoggerConfig struct {
	Level  zap.AtomicLevel
	Output []string
}

func init() {
	loggerMap = make(map[string]*zap.Logger)

	// env keeps every logger the same
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "stdout"
	}

	LogLevel := LevelInfo
	_logLevel := os.Getenv("LOG_LEVEL")
	if _logLevel != "" {
		LogLevel = map[string]zap.AtomicLevel{
			"INFO":  LevelInfo,
			"ERROR": LevelError,
			"FATAL": LevelFatal,
		}[_logLevel]
	}

	defaultConfig = &LoggerConfig{
		Level:  LogLevel,
		Output: []string{logPath},
	}
}
