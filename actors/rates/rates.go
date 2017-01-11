package rates

import (
	"time"

	"github.com/AsynkronIT/gam/actor"
	"github.com/shumkovdenis/actor/actors/group"
	"github.com/shumkovdenis/actor/manifest"
)

// var logger = actor.ActorLogger("rates")

// Change -> event.rates.change
type Change struct {
}

// Fail -> event.rates.fail
type Fail struct {
	Message string `json:"message"`
}

type ratesActor struct {
	urlAPI      string
	getInterval time.Duration
	listener    *actor.PID
	members     int
	ticker      *time.Ticker
}

func New(listener *actor.PID) actor.Actor {
	conf := manifest.Get().Config.RatesAPI
	return &ratesActor{
		urlAPI:      conf.URL,
		getInterval: conf.GetInterval,
		listener:    listener,
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
	state.ticker = time.NewTicker(state.getInterval)
	for _ = range state.ticker.C {
		state.listener.Tell(&Change{})
	}
}
