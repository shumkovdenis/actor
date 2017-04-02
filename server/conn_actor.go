package server

import (
	"go.uber.org/zap"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/logger"
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
			logger.L().Error("login failed: get session failed",
				zap.Error(err),
			)
			ctx.Respond(&LoginFailed{})
			return
		}

		sessionPID, ok := res.(*actor.PID)
		if !ok {
			ctx.Respond(&LoginFailed{res.(core.Code)})
			return
		}

		useSession := &UseSession{
			ConnPID: ctx.Self(),
		}

		future = sessionPID.RequestFuture(useSession, Timeout)
		res, err = future.Result()
		if err != nil {
			logger.L().Error("login failed: use session failed",
				zap.Error(err),
			)
			ctx.Respond(&LoginFailed{})
			return
		}

		if _, ok := res.(*UseSessionSuccess); !ok {
			ctx.Respond(&LoginFailed{res.(core.Code)})
			return
		}

		state.sessionPID = sessionPID

		ctx.Respond(&LoginSuccess{})

		ctx.SetBehavior(state.logged)
	case core.Command:
		ctx.Respond(&NotLogged{})
	}
}

func (state *connActor) logged(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Stopped:
		future := state.sessionPID.RequestFuture(&FreeSession{}, Timeout)
		res, err := future.Result()
		if err != nil {
			logger.L().Error("conn stopped failed: free session failed",
				zap.Error(err),
			)
			return
		}

		if _, ok := res.(*FreeSessionSuccess); !ok {
			logger.L().Error("free session failed")
			return
		}

		ctx.SetBehavior(state.Receive)
	case core.Event:
		ctx.Parent().Tell(msg)
	case *Login:
		ctx.Respond(&AlreadyLogged{})
	case core.Command:
		state.sessionPID.Request(msg, ctx.Sender())
	}
}
