package session

import (
	"github.com/AsynkronIT/gam/actor"
	uuid "github.com/satori/go.uuid"
	"github.com/shumkovdenis/actor/actors/account"
	"github.com/shumkovdenis/actor/actors/group"
)

// Login -> command.login
type Login struct {
}

// LoginSuccess -> event.login.success
type LoginSuccess struct {
	Client string `json:"client"`
}

// LoginFail -> event.login.fail
type LoginFail struct {
	Message string `json:"message"`
}

// Join -> command.join
type Join struct {
	Client string `mapstructure:"client"`
}

// JoinSuccess -> event.join.success
type JoinSuccess struct {
}

// JoinFail -> event.join.fail
type JoinFail struct {
	Message string `json:"message"`
}

type sessionActor struct {
	// path       string
	clientPID  *actor.PID
	accountPID *actor.PID
}

func NewActor() actor.Actor {
	return &sessionActor{}
}

func (state *sessionActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromProducer(account.NewActor)
		state.accountPID = ctx.Spawn(props)
		ctx.Become(state.started)
	}
}

func (state *sessionActor) started(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Login:
		client := uuid.NewV4().String()
		props := actor.FromProducer(group.NewActor)
		state.clientPID = actor.SpawnNamed(props, "/clients/"+client)
		state.clientPID.Tell(&group.Join{Consumer: ctx.Parent()})
		ctx.Parent().Tell(&LoginSuccess{client})
		ctx.Become(state.joined)
	case *Join:
		state.clientPID = actor.NewLocalPID("/clients/" + msg.Client)
		state.clientPID.Tell(&group.Join{Consumer: ctx.Parent()})
		ctx.Parent().Tell(&JoinSuccess{})
		ctx.Become(state.joined)
	}
}

func (state *sessionActor) joined(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case
		*account.Auth,
		*account.Balance,
		*account.Session,
		*account.Withdraw:
		state.accountPID.Request(msg, ctx.Parent())
	}
}

// func (state *sessionActor) Receive(ctx actor.Context) {
// 	switch msg := ctx.Message().(type) {
// 	case *messages.SubscribeSuccess:
// 		if msg.Contains("event.rates.change") {
// 			pid := actor.NewLocalPID("/rates")
// 			pid.Tell(&group.Join{Consumer: ctx.Self()})
// 		}
// 	case *messages.UnsubscribeSuccess:
// 		if msg.Contains("event.rates.change") {
// 		}
// 	case *Login:
// 		var id = msg.Client
// 		if len(strings.TrimSpace(id)) == 0 {
// 			state.path = "app"
// 			props := actor.FromProducer(client.NewActor)
// 			id = uuid.NewV4().String()
// 			state.clientPID = actor.SpawnNamed(props, "/clients/"+id)
// 		} else {
// 			state.path = "web"
// 			state.clientPID = actor.NewLocalPID("/clients/" + id)
// 		}
// 		state.clientPID.Request(&client.Join{Client: id}, ctx.Self())
// 	case *client.Joined:
// 		ctx.Parent().Tell(&LoginSuccess{msg.Client})
// 		switch state.path {
// 		case "app":
// 			ctx.Become(state.App)
// 		case "web":
// 			ctx.Become(state.Web)
// 			props := actor.FromProducer(account.NewActor)
// 			state.accountPID = ctx.Spawn(props)
// 		}
// 	default:
// 		//state.clientPID.Request(msg, ctx.Parent())
// 	}
// }

// func (state *sessionActor) App(ctx actor.Context) {
// switch msg := ctx.Message().(type) {
// default:
// 	state.accountPID.Request(msg, ctx.Parent())
// }
// }

// func (state *sessionActor) Web(ctx actor.Context) {
// 	switch msg := ctx.Message().(type) {
// 	default:
// 		state.accountPID.Request(msg, ctx.Parent())
// 	}
// }
