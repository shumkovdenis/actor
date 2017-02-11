package server

import "github.com/AsynkronIT/protoactor-go/actor"

type sessionStoreActor struct {
	store SessionStore
}

func newSessionStoreActor() actor.Actor {
	return &sessionStoreActor{
		store: newSessionStore(),
	}
}

func (state *sessionStoreActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case AddSession:
		session := msg.Session
		if err := state.store.Add(session); err != nil {
			ctx.Respond(err)
			return
		}
		ctx.Respond(session)
	case GetSessionByID:
		session, err := state.store.GetByID(msg.ID)
		if err != nil {
			ctx.Respond(err)
			return
		}
		ctx.Respond(session)
	}
}

type AddSession struct {
	Session *Session
}

type GetSessionByID struct {
	ID string
}
