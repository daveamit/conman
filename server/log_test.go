package main

import (
	"testing"

	"go.uber.org/zap"
)

func TestLogStringToLevel(t *testing.T) {
	if stringToLogLevel("debug") != zap.DebugLevel {
		t.Error("debug does not translate to debug")
	}
	if stringToLogLevel("info") != zap.InfoLevel {
		t.Error("info does not translate to info")
	}
	if stringToLogLevel("warn") != zap.WarnLevel {
		t.Error("warn does not translate to warn")
	}
	if stringToLogLevel("error") != zap.ErrorLevel {
		t.Error("error does not translate to error")
	}
	if stringToLogLevel("panic") != zap.PanicLevel {
		t.Error("panic does not translate to panic")
	}
	if stringToLogLevel("fatal") != zap.FatalLevel {
		t.Error("fatal does not translate to fatal")
	}
	if stringToLogLevel("invalid-unknown") != zap.DebugLevel {
		t.Error("invalid-unknown does not translate to debug")
	}
}
