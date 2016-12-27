package account

import "github.com/AsynkronIT/gam/actor"

type Fail struct {
	Message string `json:"message"`
}

type accountActor struct {
	account  string
	password string
}

func NewActor() actor.Actor {
	return &accountActor{}
}

func (state *accountActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		ctx.Become(state.started)
	}
}

func (state *accountActor) started(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Auth:
		authMsg := auth(msg)
		ctx.Respond(authMsg)
		if _, ok := authMsg.(*AuthSuccess); ok {
			ctx.Become(state.authorized)
		}
	default:
		ctx.Respond(&Fail{"Account is not authorized"})
	}
}

func (state *accountActor) authorized(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Balance:
		ctx.Respond(balance(msg))
	default:
		ctx.Respond(&Fail{"Account already authorized"})
	}
}
