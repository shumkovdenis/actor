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

type GetJackpotsList struct{}

func (*GetJackpotsList) JackpotsListMessage() {}

func (*GetJackpotsList) Command() string {
	return "command.jackpots.list"
}

type JackpotsList struct {
	Large  float64 `json:"large"`
	Medium float64 `json:"medium"`
	Small  float64 `json:"small"`
}

func (*JackpotsList) JackpotsListMessage() {}

func (*JackpotsList) Event() string {
	return "event.jackpots.list"
}

type GetJackpotsListFailed struct{}

func (*GetJackpotsListFailed) JackpotsListMessage() {}

func (*GetJackpotsListFailed) Event() string {
	return "event.jackpots.list.failed"
}
