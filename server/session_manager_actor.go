package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
)

type sessionManagerActor struct {
	mng SessionManager
}

func newSessionManagerActor() actor.Actor {
	return &sessionManagerActor{
		mng: newSessionManager(),
	}
}

func (state *sessionManagerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *CreateSession:
		session, err := state.mng.CreateSession(msg.Conf)
		if err != nil {
			ctx.Respond(CreateSessionFail(err))

			return
		}

		props := actor.FromProducer(newSessionActor).
			WithMiddleware(plugin.Use(RegistryPlugin()))

		_, err = ctx.SpawnNamed(props, session.ID)
		if err != nil {
			log.Error(err.Error())

			ctx.Respond(CreateSessionFail(err))

			return
		}

		ctx.Respond(&CreateSessionSuccess{session})
	case *UseSession:
		if err := state.mng.UseSession(msg.SessionID); err != nil {
			ctx.Respond(UseSessionFail(err))

			return
		}

		ctx.Respond(&UseSessionSuccess{})
	}
}
