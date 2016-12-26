package broker

import (
	"log"

	"github.com/AsynkronIT/gam/actor"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/mitchellh/mapstructure"
	"github.com/shumkovdenis/actor/actors/session"
	"github.com/shumkovdenis/actor/messages"
)

type subscription struct {
	Topic string `json:"topic"`
}

type brokerActor struct {
	subs       *treeset.Set
	sessionPID *actor.PID
}

func NewActor() actor.Actor {
	return &brokerActor{subs: treeset.NewWithStringComparator()}
}

func (state *brokerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromProducer(session.NewActor)
		state.sessionPID = ctx.Spawn(props)
	case *messages.Command:
		cmdData, err := processCommand(msg)
		if err != nil {
			log.Fatalf("%s data decoding error: %v\n", msg.Type, err)
		}

		switch cmdMsg := cmdData.(type) {
		case *messages.Subscribe:
			evt := state.subscribe(cmdMsg.Topic)
			ctx.Respond(evt)
		case *messages.Unsubscribe:
			evt := state.unsubscribe(cmdMsg.Topic)
			ctx.Respond(evt)
		default:
			state.sessionPID.Request(cmdData, ctx.Self())
		}
	default:
		evt := processMessage(msg)
		if evt != nil && state.subs.Contains(evt.Type) {
			ctx.Parent().Tell(evt)
		}
	}
}

func (state *brokerActor) subscribe(topic string) *messages.Event {
	state.subs.Add(topic)

	evt := &messages.Event{
		Type: "event.subscribe.success",
		Data: subscription{
			Topic: topic,
		},
	}

	return evt
}

func (state *brokerActor) unsubscribe(topic string) *messages.Event {
	state.subs.Remove(topic)

	evt := &messages.Event{
		Type: "event.unsubscribe.success",
		Data: subscription{
			Topic: topic,
		},
	}

	return evt
}

func processCommand(cmd *messages.Command) (interface{}, error) {
	var msg interface{}

	switch cmd.Type {
	case "command.subscribe":
		msg = &messages.Subscribe{}
	case "command.unsubscribe":
		msg = &messages.Unsubscribe{}
	case "command.login":
		msg = &session.Login{}
	case "command.auth":
	}

	if err := mapstructure.Decode(cmd.Data, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func processMessage(msg interface{}) *messages.Event {
	evt := &messages.Event{
		Data: msg,
	}

	switch msg.(type) {
	case *session.LoginSuccess:
		evt.Type = "event.login.success"
	}

	return evt
}
