package server

import "github.com/AsynkronIT/protoactor-go/actor"

type Ping struct {
}

func (*Ping) Command() string {
	return "command.ping"
}

type Pong struct {
}

func (*Pong) Event() string {
	return "event.pong"
}

type Update struct {
}

func (*Update) Event() string {
	return "event.update"
}

type sessionActor struct {
	roomPID *actor.PID
	connPID *actor.PID
}

func newSessionActor(roomPID *actor.PID) actor.Actor {
	return &sessionActor{
		roomPID: roomPID,
	}
}

// func (*sessionActor) Name() string {
// 	return "sessionActor"
// }

// func (state *sessionActor) Commands() []Command {
// 	return []Command{
// 		&Ping{},
// 	}
// }

func (state *sessionActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *UseSession:
		state.connPID = msg.ConnPID

		success := &UseSessionSuccess{}

		ctx.Respond(success)
	case *Ping:
		ctx.Respond(&Pong{})
	case Command:
		state.roomPID.Request(msg, ctx.Sender())
	}
}

type UseSession struct {
	ConnPID *actor.PID
}

type UseSessionSuccess struct {
}
