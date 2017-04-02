package logger

import (
	"strings"

	"go.uber.org/zap"
)

var l *zap.Logger

func init() {
	InitProduction("")
}

func InitDevelopment(file string) {
	lvl := zap.NewAtomicLevel()
	lvl.SetLevel(zap.DebugLevel)

	paths := []string{"stderr"}
	if len(strings.TrimSpace(file)) > 0 {
		paths = append(paths, strings.TrimSpace(file))
	}

	conf := zap.Config{
		Level:            lvl,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      paths,
		ErrorOutputPaths: paths,
	}

	l, _ = conf.Build()
}

func InitProduction(file string) {
	paths := []string{"stderr"}
	if len(strings.TrimSpace(file)) > 0 {
		paths = append(paths, strings.TrimSpace(file))
	}

	conf := zap.Config{
		Level:       zap.NewAtomicLevel(),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      paths,
		ErrorOutputPaths: paths,
	}

	l, _ = conf.Build()
}

func L() *zap.Logger {
	return l
}
