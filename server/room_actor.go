package server

import "github.com/AsynkronIT/protoactor-go/actor"

// type Action struct{}

// func (*Action) Command() string {
// 	return "command.action"
// }

// type ActionSuccess struct{}

// func (*ActionSuccess) Event() string {
// 	return "event.action.success"
// }

type roomActor struct {
	*Room
	// sessions []*actor.PID
}

func newRoomActor(room *Room) actor.Actor {
	return &roomActor{
		Room: room,
		// sessions: make([]*actor.PID, 0, 1),
	}
}

func (state *roomActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *JoinRoom:
		if err := state.Join(msg.Session); err != nil {
			ctx.Respond(JoinRoomFail(err))

			return
		}

		ctx.Respond(&JoinRoomSuccess{})
}
