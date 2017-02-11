package server

import "github.com/AsynkronIT/protoactor-go/actor"

type brokerActor struct {
	broker  Broker
	connPID *actor.PID
}

func newBrokerActor() actor.Actor {
	return &brokerActor{
		broker: newBroker(),
	}
}

func (state *brokerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromProducer(newConnActor)
		pid, _ := ctx.SpawnNamed(props, "conn")

		state.connPID = pid
	case *Subscribe:
		state.broker.AddTopics(msg.Topics)

		ctx.Respond(&SubscribeSuccess{msg.Topics})
	case *Unsubscribe:
		state.broker.RemoveTopics(msg.Topics)

		ctx.Respond(&UnsubscribeSuccess{msg.Topics})
	case Command:
		state.connPID.Request(msg, ctx.Sender())
	case Event:
		if state.broker.Contains(msg.Event()) {
			ctx.Parent().Tell(msg)
		}
	}
}
