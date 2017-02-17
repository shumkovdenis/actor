package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/maps/treemap"
)

type roomManagerActor struct {
	rooms *treemap.Map
}

func newRoomManagerActor() actor.Actor {
	return &roomManagerActor{
		rooms: treemap.NewWithStringComparator(),
	}
}

func (state *roomManagerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *CreateRoom:
		props := actor.FromProducer(newRoomActor)
		pid, _ := ctx.SpawnNamed(props, msg.RoomID)

		state.rooms.Put(msg.RoomID, pid)

		ctx.Respond(&Room{msg.RoomID})
	case *GetRoom:
		pid, ok := state.rooms.Get(msg.RoomID)
		if !ok {
			ctx.Respond(&RoomNotFound{})
			return
		}

		ctx.Respond(pid)
	}
}
