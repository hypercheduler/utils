package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger(name string) *zap.Logger {
	logger, found := loggerMap[name]
	if found {
		return logger
	}

	lock.Lock()
	defer lock.Unlock()
	_, check := loggerMap[name]
	if !check {
		config := &zap.Config{
			Level:         defaultConfig.Level,
			Encoding:      "json",
			OutputPaths:   defaultConfig.Output,
			InitialFields: map[string]interface{}{"label": name},
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:  "msg",
				LevelKey:    "level",
				TimeKey:     "time",
				EncodeTime:  zapcore.ISO8601TimeEncoder,
				EncodeLevel: zapcore.LowercaseLevelEncoder,
			},
		}
		logger, _ := config.Build()
		loggerMap[name] = logger
	}
	return loggerMap[name]
}
