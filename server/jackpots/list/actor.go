package list

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/shumkovdenis/club/config"
)

type listActor struct {
	sessions *hashset.Set
	ticker   *time.Ticker
}

func NewActor() actor.Actor {
	return &listActor{
		sessions: hashset.New(),
	}
}

func (state *listActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Join:
		state.sessions.Add(msg.SessionPID)
		if state.sessions.Size() == 1 {
			state.start()
		}
	case *Leave:
		state.sessions.Remove(msg.SessionPID)
		if state.sessions.Size() == 0 {
			state.stop()
		}
	case *Get:
		res := list()
		ctx.Respond(res)
	}
}

func (state *listActor) start() {
	conf := config.AccountAPI()

	state.ticker = time.NewTicker(conf.JackpotsListInterval * time.Millisecond)

	go func() {
		for _ = range state.ticker.C {
			res := list()
			state.tell(res)
		}
	}()
}

func (state *listActor) stop() {
	if state.ticker != nil {
		state.ticker.Stop()
		state.ticker = nil
	}
}

func (state *listActor) tell(msg interface{}) {
	for _, pid := range state.sessions.Values() {
		pid.(*actor.PID).Tell(msg)
	}
}

func list() Message {
	return &List{
		Large:  145624.00,
		Medium: 25812.00,
		Small:  3628.00,
	}
}
