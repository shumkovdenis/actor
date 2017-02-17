package rates

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/logger"
)

var log = logger.Get()

type Message interface {
	RatesMessage()
}

type Join struct {
	SessionPID *actor.PID
}

func (*Join) RatesMessage() {}

type Leave struct {
	SessionPID *actor.PID
}

func (*Leave) RatesMessage() {}

type Rates struct {
	Rates []*Rate `json:"rates"`
}

func (*Rates) RatesMessage() {}

func (*Rates) Event() string {
	return "event.rates"
}

type RatesFailed struct{}

func (*RatesFailed) RatesMessage() {}

func (*RatesFailed) Event() string {
	return "event.rates.failed"
}

func (*RatesFailed) Fail() string {
	return "rates_failed"
}

type Rate struct {
	Timestamp uint64  `json:"timestamp"`
	Value     float64 `json:"value"`
}
