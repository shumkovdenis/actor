package logger

import "github.com/uber-go/zap"

var l zap.Logger

func init() {
	l = zap.New(
		zap.NewTextEncoder(),
		// zap.AddCaller(),
		zap.DebugLevel,
	)
}

func Get() zap.Logger {
	return l
}
