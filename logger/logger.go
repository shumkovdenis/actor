package logger

import "github.com/uber-go/zap"

var logger zap.Logger

func init() {
	logger = zap.New(
		zap.NewTextEncoder(),
		zap.AddCaller(),
		zap.DebugLevel,
	)
}

func Get() zap.Logger {
	return logger
}
