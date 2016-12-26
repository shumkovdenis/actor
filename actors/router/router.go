package router

import (
	"log"

	"github.com/AsynkronIT/gam/actor"
	"github.com/mitchellh/mapstructure"
	"github.com/shumkovdenis/actor/actors/commands/auth"
	"github.com/shumkovdenis/actor/actors/commands/balance"
	"github.com/shumkovdenis/actor/messages"
)

type routerActor struct {
}

func NewActor() actor.Actor {
	return &routerActor{}
}

func (state *routerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.Command:
		processCommand(ctx, msg)
	case *messages.Subscribe:
		if msg.Topic == "event.balance.success" {
			props := actor.FromProducer(balance.NewActor)
			ctx.Spawn(props)
		}
	default:
		processEvent(ctx)
	}
}

func processCommand(ctx actor.Context, msg *messages.Command) {
	switch msg.Type {
	case "command.login":
	case "command.auth":
		authMsg := &auth.Auth{}
		if err := mapstructure.Decode(msg.Data, authMsg); err != nil {
			log.Fatalf("auth data decoding error: %v\n", err)
		}
		props := actor.FromProducer(auth.NewActor)
		pid := ctx.Spawn(props)
		pid.Request(authMsg, ctx.Self())
		pid.Stop()
	}
}

func processEvent(ctx actor.Context) {
	outMsg := &messages.Event{
		Data: ctx.Message(),
	}
	switch ctx.Message().(type) {
	case *auth.AuthSuccess:
		outMsg.Type = "event.auth.success"
	case *balance.BalanceSuccess:
		outMsg.Type = "event.balance.success"
	}

	if len(outMsg.Type) > 0 {
		ctx.Parent().Tell(outMsg)
	}
}
