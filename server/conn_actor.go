package server

import "github.com/AsynkronIT/protoactor-go/actor"
import "time"

type Conn interface {
	Broker() Broker
}

type connActor struct {
	broker     Broker
	sessionPID *actor.PID
}

func newConnActor() actor.Actor {
	return &connActor{
		broker: newBroker(),
	}
}

func (state *connActor) Broker() Broker {
	return state.broker
}

func (*connActor) Name() string {
	return "connActor"
}

func (*connActor) Commands() []Command {
	return []Command{
		&Subscribe{},
		&Unsubscribe{},
		&Login{},
	}
}

func (state *connActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
	case *Subscribe:
		state.broker.AddTopics(msg.Topics)

		ctx.Respond(&SubscribeSuccess{msg.Topics})
	case *Unsubscribe:
		state.broker.RemoveTopics(msg.Topics)

		ctx.Respond(&UnsubscribeSuccess{msg.Topics})
	case *Login:
		pid := actor.NewLocalPID("server/sessions")

		future := pid.RequestFuture(&UseSession{msg.SessionID}, 1*time.Second)
		res, err := future.Result()
		if err != nil {
			ctx.Respond(&LoginFail{err.Error()})

			return
		}

		if msg, ok := res.(*UseSessionFail); ok {
			ctx.Respond(&LoginFail{msg.Message})

			return
		}

		state.sessionPID = actor.NewLocalPID("server/sessions/" + msg.SessionID)

		ctx.Respond(&LoginSuccess{})
	case Command:
		state.sessionPID.Request(msg, ctx.Sender())
	case Event:
	}
}
