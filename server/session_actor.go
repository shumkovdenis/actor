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
	*Session
}

func newSessionActor(session *Session) actor.Actor {
	return &sessionActor{
		Session: session,
	}
}

func (*sessionActor) Name() string {
	return "sessionActor"
}

func (state *sessionActor) Commands() []Command {
	return []Command{
		&Ping{},
	}
}

func (state *sessionActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		// go func() {
		// 	ticker := time.Tick(5 * time.Second)
		// 	for _ = range ticker {
		// 		ctx.Parent().Tell(&Update{})
		// 	}
		// }()
	case *Ping:
		ctx.Respond(&Pong{})
	}
}
