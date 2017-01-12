package rates

import (
	"time"

	"github.com/AsynkronIT/gam/actor"
	"github.com/shumkovdenis/club/actors/group"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
)

var log = logger.Get()

// Change -> event.rates.change
type Change struct {
}

// Fail -> event.rates.fail
type Fail struct {
	Message string `json:"message"`
}

type ratesActor struct {
	listener *actor.PID
	members  int
	ticker   *time.Ticker
}

func New(listener *actor.PID) actor.Actor {
	return &ratesActor{
		listener: listener,
	}
}

func (state *ratesActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		state.listener.Request(&group.Use{
			Producer: ctx.Self(),
			Types: []interface{}{
				&Change{},
				&Fail{},
			},
		}, ctx.Self())

		ctx.Become(state.started)

		log.Info("Start rates")
	}
}

func (state *ratesActor) started(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *group.Joined:
		state.members++
		if state.members == 1 {
			go state.request()
		}
	case *group.Left:
		state.members--
		if state.members == 0 {
			state.ticker.Stop()
		}
	}
}

func (state *ratesActor) request() {
	state.ticker = time.NewTicker(config.RatesAPI().GetInterval * time.Millisecond)
	for _ = range state.ticker.C {
		state.listener.Tell(&Change{})

		log.Debug("Change rates")
	}
}
