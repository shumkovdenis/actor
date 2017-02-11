package server

import "github.com/AsynkronIT/protoactor-go/actor"

type roomStoreActor struct {
	store RoomStore
}

func newRoomStoreActor() actor.Actor {
	return &roomStoreActor{
		store: newRoomStore(),
	}
}

func (state *roomStoreActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *GetRoomByID:
		room, err := state.store.GetByID(msg.ID)
		if err != nil {
			ctx.Respond(err)
			return
		}
		ctx.Respond(room)
	case *GetRooms:
		rooms, err := state.store.GetAll()
		if err != nil {
			ctx.Respond(err)
			return
		}
		ctx.Respond(rooms)
	case *AddRoom:
		room := msg.Room
		if err := state.store.Add(room); err != nil {
			ctx.Respond(err)
			return
		}
		ctx.Respond(room)
	case *UpdateRoom:
		room := msg.Room
		if err := state.store.Update(room); err != nil {
			ctx.Respond(err)
			return
		}
		ctx.Respond(room)
	}
}

type GetRoomByID struct {
	ID string
}

type GetRooms struct{}

type AddRoom struct {
	Room *Room
}

type UpdateRoom struct {
	Room *Room
}
