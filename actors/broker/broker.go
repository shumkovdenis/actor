package broker

import (
	"log"

	"github.com/AsynkronIT/gam/actor"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/mitchellh/mapstructure"
	"github.com/shumkovdenis/actor/actors/account"
	"github.com/shumkovdenis/actor/actors/session"
	"github.com/shumkovdenis/actor/messages"
)

type subscription struct {
	Topics []string `json:"topics"`
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
			evt := state.subscribe(cmdMsg.Topics)
			ctx.Respond(evt)
		case *messages.Unsubscribe:
			evt := state.unsubscribe(cmdMsg.Topics)
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

func (state *brokerActor) subscribe(topics []string) *messages.Event {
	for _, topic := range topics {
		state.subs.Add(topic)
	}

	evt := &messages.Event{
		Type: "event.subscribe.success",
		Data: subscription{
			Topics: topics,
		},
	}

	return evt
}

func (state *brokerActor) unsubscribe(topics []string) *messages.Event {
	for _, topic := range topics {
		state.subs.Remove(topic)
	}

	evt := &messages.Event{
		Type: "event.unsubscribe.success",
		Data: subscription{
			Topics: topics,
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
	case "command.account.auth":
		msg = &account.Auth{}
	case "command.account.balance":
		msg = &account.Balance{}
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
	case *account.Fail:
		evt.Type = "event.account.fail"
	case *account.AuthSuccess:
		evt.Type = "event.account.auth.success"
	case *account.BalanceSuccess:
		evt.Type = "event.account.balance.success"
	}

	return evt
}
