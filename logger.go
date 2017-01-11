package actor

import (
	"time"

	"github.com/uber-go/zap"
)

var logger zap.Logger

func init() {
	logger = zap.New(
		zap.NewTextEncoder(zap.TextTimeFormat(time.RFC3339)),
		zap.DebugLevel,
	)
}

func ActorLogger(name string) zap.Logger {
	return logger.With(zap.String("actor", name))
}
