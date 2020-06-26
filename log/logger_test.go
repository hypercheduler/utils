package log

import (
	"testing"
)

func TestLoggerNormalUse(t *testing.T) {
	logger := GetLogger("test", "0.3")
	logger.Info("this is a test")
	logger.Error("this is an error")
	logger.Warn("this is a warning")
}
