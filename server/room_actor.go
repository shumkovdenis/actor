package server

import "github.com/AsynkronIT/protoactor-go/actor"

type Action struct{}

func (*Action) Command() string {
	return "command.action"
}

type ActionSuccess struct{}

func (*ActionSuccess) Event() string {
	return "event.action.success"
}

type roomActor struct {
	sessions []*actor.PID
}

func newRoomActor() actor.Actor {
	return &roomActor{
		sessions: make([]*actor.PID, 0, 1),
	}
}

func (state *roomActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Action:
		ctx.Self().Tell(&ActionSuccess{})
	case Event:
		for _, pid := range state.sessions {
			pid.Tell(msg)
		}
	}
}
