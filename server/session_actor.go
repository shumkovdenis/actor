package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/server/account"
	"github.com/shumkovdenis/club/server/core"
	"github.com/shumkovdenis/club/server/rates"
	"github.com/shumkovdenis/club/server/update"
)

type sessionActor struct {
	roomPID         *actor.PID
	connPID         *actor.PID
	autoUpdatePID   *actor.PID
	ratesPID        *actor.PID
	jackpotsTopsPID *actor.PID
	jackpotsListPID *actor.PID
	accountPID      *actor.PID
}

func newSessionActor(roomPID *actor.PID) actor.Actor {
	return &sessionActor{
		roomPID: roomPID,
	}
}

func (state *sessionActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.autoUpdatePID = actor.NewLocalPID("update/auto")
		state.ratesPID = actor.NewLocalPID("rates")
		state.jackpotsTopsPID = actor.NewLocalPID("jackpots/tops")
		state.jackpotsListPID = actor.NewLocalPID("jackpots/list")
	case *UseSession:
		state.connPID = msg.ConnPID

		props := actor.FromProducer(account.NewActor)
		pid, _ := ctx.SpawnNamed(props, "account")

		state.accountPID = pid

		state.autoUpdatePID.Tell(&update.Join{ctx.Self()})

		ctx.Respond(&UseSessionSuccess{})

		ctx.SetBehavior(state.used)
	case *FreeSession:
		ctx.Respond(&SessionNotUsed{})
	}
}

func (state *sessionActor) used(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case core.Event:
		state.connPID.Tell(msg)
	case *UseSession:
		ctx.Respond(&SessionAlreadyUsed{})
	case *FreeSession:
		state.accountPID.Stop()

		state.autoUpdatePID.Tell(&update.Leave{ctx.Self()})

		state.ratesPID.Tell(&rates.Leave{ctx.Self()})

		state.accountPID.Tell(&account.StopLiveJackpotsTops{})
		state.accountPID.Tell(&account.StopLiveJackpotsList{})
		state.accountPID.Tell(&account.StopLiveBalance{})

		ctx.Respond(&FreeSessionSuccess{})

		ctx.SetBehavior(state.Receive)
	case *Subscribe:
		if msg.Contains(&rates.Rates{}) {
			state.ratesPID.Tell(&rates.Join{ctx.Self()})
		}
		if msg.Contains(&account.Balance{}) {
			state.accountPID.Tell(&account.StartLiveBalance{})
		}
		if msg.Contains(&account.JackpotsTops{}) {
			state.accountPID.Tell(&account.StartLiveJackpotsTops{})
		}
		if msg.Contains(&account.JackpotsList{}) {
			state.accountPID.Tell(&account.StartLiveJackpotsList{})
		}
	case *Unsubscribe:
		if msg.Contains(&rates.Rates{}) {
			state.ratesPID.Tell(&rates.Leave{ctx.Self()})
		}
		if msg.Contains(&account.Balance{}) {
			state.accountPID.Tell(&account.StopLiveBalance{})
		}
		if msg.Contains(&account.JackpotsTops{}) {
			state.accountPID.Tell(&account.StopLiveJackpotsTops{})
		}
		if msg.Contains(&account.JackpotsList{}) {
			state.accountPID.Tell(&account.StopLiveJackpotsList{})
		}
	case account.Message:
		state.accountPID.Request(msg, ctx.Sender())
	case core.Command:
		state.roomPID.Request(msg, ctx.Sender())
	}
}
