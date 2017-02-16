package server

import "github.com/AsynkronIT/protoactor-go/actor"

const (
	SessionAlreadyUse = "session_already_use"
)

type UseSession struct {
	ConnPID *actor.PID
}

type UseSessionSuccess struct {
}

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

		props := actor.FromProducer(newAccountActor)
		pid, _ := ctx.SpawnNamed(props, "account")

		state.accountPID = pid

		ctx.Respond(&UseSessionSuccess{})

		ctx.SetBehavior(state.used)
		// case *AccountAuth, *AccountBalance, *AccountSession, *AccountWithdraw:
		// 	state.accountPID.Request(msg, ctx.Sender())
		// case *Subscribe:
		// 	state.ratesPID.Tell(&JoinRates{ctx.Self()})
		// case *Unsubscribe:
		// 	state.ratesPID.Tell(&LeaveRates{ctx.Self()})
		// case Command:
		// 	state.roomPID.Request(msg, ctx.Sender())
		// case Event:
		// 	state.connPID.Tell(msg)
	}
}

func (state *sessionActor) used(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *UseSession:
		ctx.Respond(newFail(SessionAlreadyUse))
	case *Subscribe:
		state.ratesPID.Tell(&JoinRates{ctx.Self()})
	case *Unsubscribe:
		state.ratesPID.Tell(&LeaveRates{ctx.Self()})
	case AccountCommand:
		state.accountPID.Request(msg, ctx.Sender())
	case Command:
		state.roomPID.Request(msg, ctx.Sender())
	case Event:
		state.connPID.Tell(msg)
	}
}
