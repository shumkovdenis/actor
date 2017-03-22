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

type GetJackpotsTops struct{}

func (*GetJackpotsTops) JackpotsTopsMessage() {}

func (*GetJackpotsTops) Command() string {
	return "command.jackpots.tops"
}

type JackpotsTops struct {
	Tops []Jackpot
}

func (*JackpotsTops) JackpotsTopsMessage() {}

func (*JackpotsTops) Event() string {
	return "event.jackpots.tops"
}

type GetJackpotsTopsFailed struct{}

func (*GetJackpotsTopsFailed) JackpotsTopsMessage() {}

func (*GetJackpotsTopsFailed) Event() string {
	return "event.jackpots.tops.failed"
}

type Jackpot struct {
	Account string  `json:"account"`
	Win     float64 `json:"win"`
	Date    int64   `json:"date"`
}
