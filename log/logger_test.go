package log

import (
	"testing"
)

func TestLoggerNormalUse(t *testing.T) {
	logger := GetLogger("test")
	logger.Info("this is a test")
	logger.Error("this is an error")
	logger.Warn("this is a warning")
}
