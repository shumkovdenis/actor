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
		}

		state.sessions.Put(id, ctx.Sender())
		state.connections.Put(pid, ctx.Sender())

		ctx.Respond(&CreateSessionSuccess{id})
	case *UseSession:
		if _, ok := state.sessions.Get(msg.SessionID); !ok {
			ctx.Respond(&UseSessionFail{"Session not found"})

			return
		}

		ctx.Respond(&UseSessionSuccess{})
	case Event:
		if conn, ok := state.connections.Get(ctx.Sender()); ok {
			conn.(*actor.PID).Tell(msg)
		}
	}
}
