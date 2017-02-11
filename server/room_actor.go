package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/sets/hashset"
)

// type Action struct{}

// func (*Action) Command() string {
// 	return "command.action"
// }

// type ActionSuccess struct{}

// func (*ActionSuccess) Event() string {
// 	return "event.action.success"
// }

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
			ctx.Respond(&Fail{Code: RoomFull})
			return
		}

		state.sessions.Add(msg.SessionPID)

		success := &JoinRoomSuccess{}

		ctx.Respond(success)
	}
}

type JoinRoom struct {
	SessionPID *actor.PID
}

type JoinRoomSuccess struct {
}
