package logger

import "go.uber.org/zap"

var l *zap.Logger

func init() {
	l, _ = zap.NewProduction()
}

func Get() *zap.Logger {
	return l
}
