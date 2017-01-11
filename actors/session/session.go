package session

import (
	"github.com/AsynkronIT/gam/actor"
	uuid "github.com/satori/go.uuid"
	"github.com/shumkovdenis/club/actors/account"
	"github.com/shumkovdenis/club/actors/group"
	"github.com/shumkovdenis/club/messages"
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
	client     string
	clientPID  *actor.PID
	accountPID *actor.PID
}

func NewActor() actor.Actor {
	return &sessionActor{}
}

func (state *sessionActor) Receive(ctx actor.Context) {
	state.subscription(ctx)

	switch msg := ctx.Message().(type) {
	case *actor.Started:
		// pid := actor.NewLocalPID("/update")
		// pid.Tell(&group.Join{Consumer: ctx.Parent()})
	case *Login:
		state.client = uuid.NewV4().String()

		props := actor.FromProducer(group.NewActor)
		state.clientPID = actor.SpawnNamed(props, "/clients/"+state.client)

		ctx.Become(state.use)

		state.clientPID.Request(&group.Use{Producer: ctx.Self()}, ctx.Self())
	case *Join:
		state.clientPID = actor.NewLocalPID("/clients/" + msg.Client)

		ctx.Become(state.use)

		state.clientPID.Request(&group.Use{Producer: ctx.Self()}, ctx.Self())
	}
}

func (state *sessionActor) use(ctx actor.Context) {
	state.subscription(ctx)

	switch ctx.Message().(type) {
	case *group.Used:
		state.clientPID.Tell(&group.Join{Consumer: ctx.Parent()})

		props := actor.FromProducer(account.NewActor)
		state.accountPID = ctx.Spawn(props)

		ctx.Become(state.used)

		if len(state.client) > 0 {
			ctx.Parent().Tell(&LoginSuccess{state.client})
		} else {
			ctx.Parent().Tell(&JoinSuccess{})
		}
	}
}

func (state *sessionActor) used(ctx actor.Context) {
	state.subscription(ctx)

	switch msg := ctx.Message().(type) {
	case
		*account.Auth,
		*account.Balance,
		*account.Session,
		*account.Withdraw:
		state.accountPID.Request(msg, ctx.Parent())
	}
}

func (state *sessionActor) subscription(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.SubscribeSuccess:
		if msg.Contains("event.rates.change") {
			pid := actor.NewLocalPID("/rates")
			pid.Tell(&group.Join{Consumer: ctx.Parent()})
		}
	case *messages.UnsubscribeSuccess:
		if msg.Contains("event.rates.change") {
			pid := actor.NewLocalPID("/rates")
			pid.Tell(&group.Leave{Consumer: ctx.Parent()})
		}
	}
}
