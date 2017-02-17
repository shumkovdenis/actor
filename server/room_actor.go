package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/sets/hashset"
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
			ctx.Respond(&RoomFull{})
			return
		}

		state.sessions.Add(msg.SessionPID)

		ctx.Respond(&JoinedRoom{})
	case *LeaveRoom:
		state.sessions.Remove(msg.SessionPID)

		ctx.Respond(&LeftRoom{})
	}
}
