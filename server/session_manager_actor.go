package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/maps/treemap"
	uuid "github.com/satori/go.uuid"
)

type sessionManagerActor struct {
	sessions    *treemap.Map
	connections *hashmap.Map
}

func newSessionManagerActor() actor.Actor {
	return &sessionManagerActor{
		sessions:    treemap.NewWithStringComparator(),
		connections: hashmap.New(),
	}
}

func (state *sessionManagerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *CreateSession:
		id := uuid.NewV4().String()

		props := actor.FromProducer(newSessionActor).
			WithMiddleware(plugin.Use(RegistryPlugin()))
		pid, err := ctx.SpawnNamed(props, id)
		if err != nil {
			ctx.Respond(&CreateSessionFail{err.Error()})

			return
		}

		state.sessions.Put(id, pid)

		ctx.Respond(&CreateSessionSuccess{id})
	case *UseSession:
		pid, ok := state.sessions.Get(msg.SessionID)
		if !ok {
			ctx.Respond(&UseSessionFail{"Session not found"})

			return
		}

		state.connections.Put(ctx.Sender(), pid)

		ctx.Respond(&UseSessionSuccess{})
	}
}
