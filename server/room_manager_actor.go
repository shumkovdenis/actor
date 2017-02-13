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
		id := msg.RoomID

		props := actor.FromProducer(newRoomActor)
		pid, _ := ctx.SpawnNamed(props, id)

		state.rooms.Put(id, pid)

		success := &CreateRoomSuccess{
			Room: &Room{
				ID: id,
			},
		}

		ctx.Respond(success)
	case *GetRoom:
		pid, ok := state.rooms.Get(msg.RoomID)
		if !ok {
			err := newErr(ErrRoomNotFound).LogErr()
			err = newErr(ErrGetRoom).Wrap(err).LogErr()
			ctx.Respond(err)
			return
		}

		success := &GetRoomSuccess{
			RoomPID: pid.(*actor.PID),
		}

		ctx.Respond(success)
	}
}

type CreateRoom struct {
	RoomID string
}

type CreateRoomSuccess struct {
	Room *Room
}

type GetRoom struct {
	RoomID string
}

type GetRoomSuccess struct {
	RoomPID *actor.PID
}
