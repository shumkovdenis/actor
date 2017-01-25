package session

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/actors/account"
	"github.com/shumkovdenis/club/actors/app/update"
	"github.com/shumkovdenis/club/logger"
	"github.com/uber-go/zap"
)

var log = logger.Get()

type Session interface {
	ID() string
}

type sessionActor struct {
	id         string
	updatePID  *actor.PID
	accountPID *actor.PID
}

func New(id string) actor.Actor {
	return &sessionActor{id: id}
}

func (state *sessionActor) ID() string {
	return state.id
}

func (state *sessionActor) Receive(ctx actor.Context) {
	if ok := state.update(ctx); ok {
		return
	}

	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.updatePID = actor.NewLocalPID("update")

		log.Info("Session actor started",
			zap.String("id", state.id),
		)
	case *Login:
		props := actor.FromProducer(account.NewActor)
		state.accountPID = ctx.Spawn(props)

		ctx.SetBehavior(state.used)
	default:
		if _, ok := msg.(actor.SystemMessage); !ok {
			// ctx.Respond(&Fail{"The session is not logged"})
		}
	}
}

func (state *sessionActor) used(ctx actor.Context) {
	if ok := state.update(ctx); ok {
		return
	}

	switch msg := ctx.Message().(type) {
	case
		*account.Auth,
		*account.Balance,
		*account.Session,
		*account.Withdraw:
		state.accountPID.Request(msg, ctx.Sender())
	}
}

func (state *sessionActor) update(ctx actor.Context) bool {
	switch msg := ctx.Message().(type) {
	case
		*update.Check,
		*update.Download,
		*update.Install:
		state.updatePID.Request(msg, ctx.Sender())

		return true
	}

	return false
}

/*
type SessionActor struct {
	client     string
	clientPID  *actor.PID
	accountPID *actor.PID
}

func NewActor() actor.Actor {
	return &SessionActor{}
}

func (state *SessionActor) Receive(ctx actor.Context) {
	state.subscription(ctx)
	state.update(ctx)

	switch msg := ctx.Message().(type) {
	case *actor.Started:
		pid := actor.NewLocalPID("update/auto")
		pid.Tell(&group.Join{Consumer: ctx.Parent()})
	case *Login:
		state.client = uuid.NewV4().String()

		props := actor.FromProducer(group.New)
		state.clientPID, _ = actor.SpawnNamed(props, "/clients/"+state.client)

		ctx.SetBehavior(state.use)

		state.clientPID.Request(&group.Use{Producer: ctx.Self()}, ctx.Self())
	case *Join:
		state.clientPID = actor.NewLocalPID("/clients/" + msg.Client)

		ctx.SetBehavior(state.use)

		state.clientPID.Request(&group.Use{Producer: ctx.Self()}, ctx.Self())
	}
}

func (state *SessionActor) use(ctx actor.Context) {
	state.subscription(ctx)
	state.update(ctx)

	switch ctx.Message().(type) {
	case *group.Used:
		state.clientPID.Tell(&group.Join{Consumer: ctx.Parent()})

		props := actor.FromProducer(account.NewActor)
		state.accountPID = ctx.Spawn(props)

		ctx.SetBehavior(state.used)

		if len(state.client) > 0 {
			ctx.Parent().Tell(&LoginSuccess{state.client})
		} else {
			ctx.Parent().Tell(&JoinSuccess{})
		}
	}
}

func (state *SessionActor) used(ctx actor.Context) {
	state.subscription(ctx)
	state.update(ctx)

	switch msg := ctx.Message().(type) {
	case
		*account.Auth,
		*account.Balance,
		*account.Session,
		*account.Withdraw:
		state.accountPID.Request(msg, ctx.Parent())
	}
}

func (state *SessionActor) subscription(ctx actor.Context) {
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

func (state *SessionActor) update(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case
		*update.Check,
		*update.Download,
		*update.Install:
		actor.NewLocalPID("/update").Request(msg, ctx.Parent())
	}
}
*/
