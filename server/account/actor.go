package account

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/server/core"
)

type accountActor struct {
	username string
	password string
}

func NewActor() actor.Actor {
	return &accountActor{}
}

func (state *accountActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Authorize:
		res := authorize(msg.Username, msg.Password)
		ctx.Respond(res)
		if !core.IsFail(res) {
			state.username = msg.Username
			state.password = msg.Password
			ctx.SetBehavior(state.authorized)
		}
	case Message:
		ctx.Respond(&NotAuthorized{})
	}
}

func (state *accountActor) authorized(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Authorize:
		ctx.Respond(&AlreadyAuthorized{})
	case *GetBalance:
		res := getBalance(state.username, state.password)
		ctx.Respond(res)
	case *GetGameSession:
		res := getGameSession(state.username, state.password, msg.GameID)
		ctx.Respond(res)
	case *Withdraw:
		res := withdraw(state.username, state.password)
		ctx.Respond(res)
		if !core.IsFail(res) {
			ctx.SetBehavior(state.Receive)
		}
	}
}
