package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func stringToLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	}

	return zap.DebugLevel
}

var atom = zap.NewAtomicLevel()
var log *zap.Logger

func setLog(level zapcore.Level, encoding string) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.StacktraceKey = "stack"
	cfg.EncoderConfig.CallerKey = "line"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg.Encoding = encoding
	cfg.Level = atom
	if encoding == "console" {
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, _ := cfg.Build(zap.AddCaller())

	// By Default log everything.
	atom.SetLevel(level)

	if log != nil {
		log.Sync()
	}

	log = logger // .Named(theService.Name) // .With(state)

	// if theService != nil {
	// 	log = log.Named(theService.Name)
	// }
}

func init() {
	setLog(stringToLogLevel("debug"), "console")
}
