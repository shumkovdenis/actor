package server

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type connActor struct {
	sessionPID *actor.PID
}

func newConnActor() actor.Actor {
	return &connActor{}
}

func (*connActor) Name() string {
	return "connActor"
}

func (*connActor) Commands() []Command {
	return []Command{
		&Login{},
	}
}

func (state *connActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *Login:
		pid := actor.NewLocalPID("server/sessions")

		future := pid.RequestFuture(&UseSession{msg.SessionID}, 1*time.Second)
		res, err := future.Result()
		if err != nil {
			ctx.Respond(&LoginFail{err.Error()})

			return
		}

		if msg, ok := res.(*UseSessionFail); ok {
			ctx.Respond(&LoginFail{msg.Message})

			return
		}

		state.sessionPID = actor.NewLocalPID("server/sessions/" + msg.SessionID)

		ctx.Respond(&LoginSuccess{})
	case Command:
		state.sessionPID.Request(msg, ctx.Sender())
	}
}
