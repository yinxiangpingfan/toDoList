package logger

import (
	"testing"
)

func TestLoggerInit(t *testing.T) {
	// LoggerInit(path )
	Logger.Errorf("2")
	Logger.Debugf("1")
	Logger.Panicf("3")
}
