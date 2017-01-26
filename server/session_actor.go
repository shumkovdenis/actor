package server

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Ping struct {
}

type Pong struct {
}

type Update struct {
}

type sessionActor struct {
}

func newSessionActor() actor.Actor {
	return &sessionActor{}
}

func (state *sessionActor) Command(typ string) interface{} {
	switch typ {
	case "command.ping":
		return &Ping{}
	}
	return nil
}

func (state *sessionActor) Event(msg interface{}) string {
	switch msg.(type) {
	case *Pong:
		return "event.pong"
	case *Update:
		return "event.update"
	}
	return ""
}

func (state *sessionActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		go func() {
			ticker := time.Tick(5 * time.Second)
			for _ = range ticker {
				ctx.Parent().Tell(&Update{})
			}
		}()
	case *Ping:
		ctx.Respond(&Pong{})
	}
}
