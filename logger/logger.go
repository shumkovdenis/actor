package logger

import "go.uber.org/zap"

var l *zap.Logger

func init() {
	l, _ = zap.NewDevelopment()
}

func Get() *zap.Logger {
	return l
}

func File(path string) *zap.Logger {
	w, _, _ := zap.Open(path)
	return l.WithOptions(zap.ErrorOutput(w))
}
