package client

import "github.com/AsynkronIT/gam/actor"

type Join struct {
	Client string
}

type Joined struct {
	Client string
}

type clientActor struct {
	sessions []*actor.PID
}

func NewActor() actor.Actor {
	return &clientActor{
		sessions: []*actor.PID{},
	}
}

func (state *clientActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Join:
		state.sessions = append(state.sessions, ctx.Sender())
		ctx.Respond(&Joined{msg.Client})
	}
}
