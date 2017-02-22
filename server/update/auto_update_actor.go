package update

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/shumkovdenis/club/config"
)

type autoUpdateActor struct {
	sessions *hashset.Set
}

func newAutoUpdateActor() actor.Actor {
	return &autoUpdateActor{
		sessions: hashset.New(),
	}
}

func (state *autoUpdateActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		conf := config.UpdateServer()
		if conf.AutoUpdate {
			state.loop(ctx)
		}
	case *Join:
		state.sessions.Add(msg.SessionPID)
	case *Leave:
		state.sessions.Remove(msg.SessionPID)
	case *Available:
		ctx.Parent().Request(&Download{}, ctx.Self())
	case *No:
		state.loop(ctx)
	case *Complete:
		ctx.Parent().Request(&Install{}, ctx.Self())
	case *Ready:
		state.tell(msg)
		state.loop(ctx)
	}
}

func (state *autoUpdateActor) tell(msg interface{}) {
	for _, pid := range state.sessions.Values() {
		pid.(*actor.PID).Tell(msg)
	}
}

func (state *autoUpdateActor) loop(ctx actor.Context) {
	conf := config.UpdateServer()

	go func() {
		time.Sleep(conf.CheckInterval * time.Millisecond)

		ctx.Parent().Request(&Check{}, ctx.Self())
	}()
}
