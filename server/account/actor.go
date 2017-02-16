package account

import "github.com/AsynkronIT/protoactor-go/actor"

type accountActor struct {
}

func NewActor() actor.Actor {
	return &accountActor{}
}

func (state *accountActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Authorize:
	case Incoming:
		ctx.Respond(&NotAuthorized{})
	}
}

func (state *accountActor) authorized(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Authorize:
		ctx.Respond(&AlreadyAuthorized{})
	}
}
