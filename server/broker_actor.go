package server

import "github.com/AsynkronIT/protoactor-go/actor"

type brokerActor struct {
	brk     Broker
	connPID *actor.PID
}

func newBrokerActor() actor.Actor {
	return &brokerActor{
		brk: newBroker(),
	}
}

func (*brokerActor) Name() string {
	return "brokerActor"
}

func (*brokerActor) Commands() []Command {
	return []Command{
		&Subscribe{},
		&Unsubscribe{},
	}
}

func (state *brokerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromProducer(newConnActor)
		pid, err := ctx.SpawnNamed(props, "conn")
		if err != nil {
		}

		state.connPID = pid
	case *Subscribe:
		state.brk.AddTopics(msg.Topics)

		ctx.Respond(&SubscribeSuccess{msg.Topics})
	case *Unsubscribe:
		state.brk.RemoveTopics(msg.Topics)

		ctx.Respond(&UnsubscribeSuccess{msg.Topics})
	case Command:
		// state.connPID.Request(msg, ctx.Sender())
		ctx.Respond(&Fail{"No proccess"})
		// case Event:
		// 	if state.brk.Contains(msg.Event()) {
		// 		ctx.Parent().Tell(msg)
		// 	}
	}
}
