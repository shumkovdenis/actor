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
			log.Error("get session fail")
			ctx.Respond(&LoginFail{newError(ErrLogin)})
			return
		}

		if err, ok := res.(*Error); ok {
			fail := &LoginFail{newErrorWrap(ErrGetSession, err)}
			ctx.Respond(fail)
			return
		}

		sessionPID := res.(*GetSessionSuccess).SessionPID

		useSession := &UseSession{
			ConnPID: ctx.Self(),
		}

		future = sessionPID.RequestFuture(useSession, 1*time.Second)
		res, err = future.Result()
		if err != nil {
			fail := &LoginFail{newErrorWrap(ErrUseSession, err)}
			ctx.Respond(fail)
			return
		}

		if err, ok := res.(*Error); ok {
			fail := &LoginFail{newErrorWrap(ErrUseSession, err)}
			ctx.Respond(fail)
			return
		}

		state.sessionPID = sessionPID

		ctx.Respond(&LoginSuccess{})
	case Command:
		state.sessionPID.Request(msg, ctx.Sender())
	case Event:
		ctx.Parent().Tell(msg)
	}
}
