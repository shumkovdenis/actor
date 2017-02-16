package server

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type connActor struct {
	sessionManagerPID *actor.PID
	sessionPID        *actor.PID
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
		state.sessionManagerPID = actor.NewLocalPID("sessions")
	case *Login:
		getSession := &GetSession{
			SessionID: msg.SessionID,
		}

		future := state.sessionManagerPID.RequestFuture(getSession, 1*time.Second)
		res, err := future.Result()
		if err != nil {
			// err := newErr(ErrLogin).Error(err).LogErr()
			// ctx.Respond(&LoginFail{err})
			return
		}

		// if err, ok := res.(*Err); ok {
		// err := newErr(ErrLogin).Wrap(err).LogErr()
		// ctx.Respond(&LoginFail{err})
		// 	return
		// }

		sessionPID := res.(*GetSessionSuccess).SessionPID

		useSession := &UseSession{
			ConnPID: ctx.Self(),
		}

		future = sessionPID.RequestFuture(useSession, 1*time.Second)
		res, err = future.Result()
		if err != nil {
			// err := newErr(ErrLogin).Error(err).LogErr()
			// ctx.Respond(&LoginFail{err})
			return
		}

		// if err, ok := res.(*Err); ok {
		// err := newErr(ErrLogin).Wrap(err).LogErr()
		// ctx.Respond(&LoginFail{err})
		// 	return
		// }

		state.sessionPID = sessionPID

		ctx.Respond(&LoginSuccess{})
	case Command:
		state.sessionPID.Request(msg, ctx.Sender())
	case Event:
		ctx.Parent().Tell(msg)
	}
}
