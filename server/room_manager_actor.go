package server

import "github.com/AsynkronIT/protoactor-go/actor"

type roomManagerActor struct {
	mng RoomManager
}

func newRoomManagerActor() actor.Actor {
	return &roomManagerActor{
		mng: newRoomManager(),
	}
}

func (state *roomManagerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *CreateRoom:
		room, err := state.mng.Create(msg.Conf)
		if err != nil {
			ctx.Respond(CreateRoomFail(err))

			return
		}

		// props := actor.FromProducer(newRoomActor)
		// _, err = ctx.SpawnNamed(props, room.ID)
		// if err != nil {
		// 	log.Error(err.Error())

		// 	ctx.Respond(CreateRoomFail(err))

		// 	return
		// }

		ctx.Respond(&CreateRoomSuccess{room})
	case *JoinRoom:
		if err := state.mng.JoinRoom(msg.RoomID); err != nil {
			ctx.Respond(JoinRoomFail(err))

			return
		}

		ctx.Respond(&JoinRoomSuccess{})
	case *LeaveRoom:
		if err := state.mng.LeaveRoom(msg.RoomID); err != nil {
			ctx.Respond(LeaveRoomFail(err))

			return
		}

		ctx.Respond(&LeaveRoomSuccess{})
	}
}
