package account

import "github.com/AsynkronIT/gam/actor"

type AccountFail struct {
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
	switch msg := ctx.Message().(type) {
	case *Auth:
		authMsg := auth(msg)
		ctx.Respond(authMsg)
		if _, ok := authMsg.(*AuthSuccess); ok {
			ctx.Become(state.authorized)
		}
	default:
		ctx.Respond(&AccountFail{"Account is not authorized"})
	}
}

func (state *accountActor) authorized(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Balance:
		balance(msg)
	default:
		ctx.Respond(&AccountFail{"Account already authorized"})
	}
}
