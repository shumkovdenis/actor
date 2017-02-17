package server

import (
	"go.uber.org/zap"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/server/core"
)

type connActor struct {
	sessionManagerPID *actor.PID
	sessionPID        *actor.PID
}

func newConnActor() actor.Actor {
	return &connActor{}
}

func (state *connActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.sessionManagerPID = actor.NewLocalPID("sessions")
	case *Login:
		getSession := &GetSession{
			SessionID: msg.SessionID,
		}

		future := state.sessionManagerPID.RequestFuture(getSession, Timeout)
		res, err := future.Result()
		if err != nil {
			log.Error("login failed: get session failed",
				zap.Error(err),
			)
			ctx.Respond(&LoginFailed{})
			return
		}

		sessionPID, ok := res.(*actor.PID)
		if !ok {
			ctx.Respond(res)
			return
		}

		useSession := &UseSession{
			ConnPID: ctx.Self(),
		}

		future = sessionPID.RequestFuture(useSession, Timeout)
		res, err = future.Result()
		if err != nil {
			log.Error("login failed: use session failed",
				zap.Error(err),
			)
			ctx.Respond(&LoginFailed{})
			return
		}

		if _, ok := res.(*UseSessionSuccess); !ok {
			ctx.Respond(res)
			return
		}

		state.sessionPID = sessionPID

		ctx.Respond(&LoginSuccess{})
	case core.Command:
		state.sessionPID.Request(msg, ctx.Sender())
	case core.Event:
		ctx.Parent().Tell(msg)
	}
}
