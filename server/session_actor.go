package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/server/account"
	"github.com/shumkovdenis/club/server/core"
	"github.com/shumkovdenis/club/server/rates"
)

type sessionActor struct {
	roomPID    *actor.PID
	connPID    *actor.PID
	ratesPID   *actor.PID
	accountPID *actor.PID
}

func newSessionActor(roomPID *actor.PID) actor.Actor {
	return &sessionActor{
		roomPID: roomPID,
	}
}

func (state *sessionActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.ratesPID = actor.NewLocalPID("rates")
	case *UseSession:
		state.connPID = msg.ConnPID

		props := actor.FromProducer(account.NewActor)
		pid, _ := ctx.SpawnNamed(props, "account")

		state.accountPID = pid

		ctx.Respond(&UseSessionSuccess{})

		ctx.SetBehavior(state.used)
	}
}

func (state *sessionActor) used(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *UseSession:
		ctx.Respond(&SessionAlreadyUse{})
	case *Subscribe:
		state.ratesPID.Tell(&rates.Join{ctx.Self()})
	case *Unsubscribe:
		state.ratesPID.Tell(&rates.Leave{ctx.Self()})
	case account.Message:
		state.accountPID.Request(msg, ctx.Sender())
	case core.Command:
		state.roomPID.Request(msg, ctx.Sender())
	case core.Event:
		state.connPID.Tell(msg)
	}
}
