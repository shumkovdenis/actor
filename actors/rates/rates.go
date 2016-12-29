package rates

import (
	"time"

	"github.com/AsynkronIT/gam/actor"
	"github.com/shumkovdenis/actor/actors/group"
	"github.com/shumkovdenis/actor/config"
)

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

func NewActor(listener *actor.PID) actor.Actor {
	return &ratesActor{listener: listener}
}

func (state *ratesActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		state.listener.Tell(&group.Use{
			Producer:  ctx.Self(),
			Validator: validator,
		})
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
	ratesAPI := config.Conf.RatesAPI
	state.ticker = time.NewTicker(time.Duration(ratesAPI.Timeout) * time.Millisecond)
	for _ = range state.ticker.C {
		state.listener.Tell(&Change{})
	}
}

func validator(msg interface{}) bool {
	switch msg.(type) {
	case *Change, *Fail:
		return true
	}
	return false
}
