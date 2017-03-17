package tops

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/logger"
)

var log = logger.Get()

type Message interface {
	JackpotsTopsMessage()
}

type Join struct {
	SessionPID *actor.PID
}

func (*Join) JackpotsTopsMessage() {}

type Leave struct {
	SessionPID *actor.PID
}

func (*Leave) JackpotsTopsMessage() {}

type Get struct{}

func (*Get) JackpotsTopsMessage() {}

func (*Get) Command() string {
	return "command.jackpots.tops"
}

type Tops struct {
	Tops []Jackpot
}

func (*Tops) JackpotsTopsMessage() {}

func (*Tops) Event() string {
	return "event.jackpots.tops"
}

type GetFailed struct{}

func (*GetFailed) JackpotsTopsMessage() {}

func (*GetFailed) Event() string {
	return "event.jackpots.tops.failed"
}

type Jackpot struct {
	Account string  `json:"account"`
	Win     float64 `json:"win"`
	Date    int64   `json:"date"`
}
