package server

import "github.com/AsynkronIT/protoactor-go/actor"

type sessionActor struct {
	roomPID    *actor.PID
	connPID    *actor.PID
	accountPID *actor.PID
	ratesPID   *actor.PID
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

		success := &UseSessionSuccess{}

		ctx.Respond(success)

		props := actor.FromProducer(newAccountActor)
		pid, _ := ctx.SpawnNamed(props, "account")

		state.accountPID = pid
	case *AccountAuth, *AccountBalance, *AccountSession, *AccountWithdraw:
		state.accountPID.Request(msg, ctx.Sender())
	case *Subscribe:
		state.ratesPID.Tell(&JoinRates{ctx.Self()})
	case *Unsubscribe:
		state.ratesPID.Tell(&LeaveRates{ctx.Self()})
	case Command:
		state.roomPID.Request(msg, ctx.Sender())
	case Event:
		state.connPID.Tell(msg)
	}
}

type UseSession struct {
	ConnPID *actor.PID
}

type UseSessionSuccess struct {
}
