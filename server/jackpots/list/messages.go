package list

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/logger"
)

var log = logger.Get()

type Message interface {
	JackpotsListMessage()
}

type Join struct {
	SessionPID *actor.PID
}

func (*Join) JackpotsListMessage() {}

type Leave struct {
	SessionPID *actor.PID
}

func (*Leave) JackpotsListMessage() {}

type Get struct{}

func (*Get) JackpotsListMessage() {}

func (*Get) Command() string {
	return "command.jackpots.list"
}

type List struct {
	Large  float64 `json:"large"`
	Medium float64 `json:"medium"`
	Small  float64 `json:"small"`
}

func (*List) JackpotsListMessage() {}

func (*List) Event() string {
	return "event.jackpots.list"
}

type GetFailed struct{}

func (*GetFailed) JackpotsListMessage() {}

func (*GetFailed) Event() string {
	return "event.jackpots.list.failed"
}
