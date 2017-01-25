package broker

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
	"github.com/shumkovdenis/club/actors/account"
	"github.com/shumkovdenis/club/actors/app/update"
	"github.com/shumkovdenis/club/actors/rates"
	"github.com/shumkovdenis/club/actors/session"
	"github.com/shumkovdenis/club/logger"
	"github.com/shumkovdenis/club/messages"
)

var log = logger.Get()

type Broker interface {
	Subs() *treeset.Set
}

type brokerActor struct {
	subs       *treeset.Set
	sessionPID *actor.PID
}

func New() actor.Actor {
	subs := treeset.NewWithStringComparator()
	subs.Add("event.subscribe.success", "event.unsubscribe.success")
	return &brokerActor{
		subs: subs,
	}
}

// func (b *brokerActor) Subs() *treeset.Set {
// 	return b.subs
// }

func (state *brokerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		id := uuid.NewV4().String()
		props := actor.FromInstance(session.New(id))
		state.sessionPID, _ = actor.SpawnNamed(props, "sessions/"+id)

		log.Info("Broker actor started")
	case *messages.Command:
		message, err := processCommand(msg)
		if err != nil {
			log.Fatal(fmt.Sprintf("%s data decoding error: %v\n", msg.Type, err))
		}

		if m := state.subscription(message); m != nil {
			// state.sessionPID.Tell(m)
			ctx.Self().Tell(m)
		} else {
			state.sessionPID.Request(message, ctx.Self())
		}
	default:
		evt := processMessage(msg)
		if evt != nil && state.subs.Contains(evt.Type) {
			ctx.Parent().Tell(evt)
		}
	}
}

func (state *brokerActor) subscription(message interface{}) interface{} {
	switch msg := message.(type) {
	case *messages.Subscribe:
		for _, topic := range msg.Topics {
			state.subs.Add(topic)
		}
		return &messages.SubscribeSuccess{msg.Topics}
	case *messages.Unsubscribe:
		for _, topic := range msg.Topics {
			state.subs.Remove(topic)
		}
		return &messages.UnsubscribeSuccess{msg.Topics}
	}
	return nil
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
	// case "command.join":
	// 	msg = &session.Join{}
	case "command.update.check":
		msg = &update.Check{}
	case "command.update.download":
		msg = &update.Download{}
	case "command.update.install":
		msg = &update.Install{}
	case "command.account.auth":
		msg = &account.Auth{}
	case "command.account.balance":
		msg = &account.Balance{}
	case "command.account.session":
		msg = &account.Session{}
	case "command.account.withdraw":
		msg = &account.Withdraw{}
	default:
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
	case *messages.SubscribeSuccess:
		evt.Type = "event.subscribe.success"
	case *messages.UnsubscribeSuccess:
		evt.Type = "event.unsubscribe.success"
	case *session.LoginSuccess:
		evt.Type = "event.login.success"
	case *session.LoginFail:
		evt.Type = "event.login.fail"
	// case *session.JoinSuccess:
	// 	evt.Type = "event.join.success"
	// case *session.JoinFail:
	// 	evt.Type = "event.join.fail"
	case *update.No:
		evt.Type = "event.update.no"
	case *update.Available:
		evt.Type = "event.update.available"
	case *update.DownloadProgress:
		evt.Type = "event.update.download.progress"
	case *update.DownloadComplete:
		evt.Type = "event.update.download.complete"
	case *update.InstallComplete:
		evt.Type = "event.update.install.complete"
	case *update.InstallRestart:
		evt.Type = "event.update.install.restart"
	case *update.Fail:
		evt.Type = "event.update.fail"
	case *account.Fail:
		evt.Type = "event.account.fail"
	case *account.AuthSuccess:
		evt.Type = "event.account.auth.success"
	case *account.AuthFail:
		evt.Type = "event.account.auth.fail"
	case *account.BalanceSuccess:
		evt.Type = "event.account.balance.success"
	case *account.BalanceFail:
		evt.Type = "event.account.balance.fail"
	case *account.SessionSuccess:
		evt.Type = "event.account.session.success"
	case *account.SessionFail:
		evt.Type = "event.account.session.fail"
	case *account.WithdrawSuccess:
		evt.Type = "event.account.withdraw.success"
	case *account.WithdrawFail:
		evt.Type = "event.account.withdraw.fail"
	case *rates.Change:
		evt.Type = "event.rates.change"
	case *rates.Fail:
		evt.Type = "event.rates.fail"
	}

	return evt
}
