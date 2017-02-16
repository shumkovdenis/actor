package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/sets/hashset"
)

const (
	RoomFull = "room_full"
)

type roomActor struct {
	sessions *hashset.Set
}

func newRoomActor() actor.Actor {
	return &roomActor{
		sessions: hashset.New(),
	}
}

func (state *roomActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *JoinRoom:
		if state.sessions.Size() == 2 {
			ctx.Respond(newFail(RoomFull))
			return
		}

		state.sessions.Add(msg.SessionPID)

		ctx.Respond(&JoinRoomSuccess{})
	}
}

type JoinRoom struct {
	SessionPID *actor.PID
}

type JoinRoomSuccess struct {
}
