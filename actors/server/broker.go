package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/messages"
)

// Subscribe -> command.subscribe
type Subscribe struct {
	Topics []string `mapstructure:"topics"`
}

// SubscribeSuccess -> event.subscribe.success
type SubscribeSuccess struct {
	Topics []string `json:"topics"`
}

func (s *SubscribeSuccess) Contains(topic string) bool {
	for _, t := range s.Topics {
		if topic == t {
			return true
		}
	}
	return false
}

// Unsubscribe -> command.unsubscribe
type Unsubscribe struct {
	Topics []string `mapstructure:"topics"`
}

// UnsubscribeSuccess -> event.unsubscribe.success
type UnsubscribeSuccess struct {
	Topics []string `json:"topics"`
}

func (s *UnsubscribeSuccess) Contains(topic string) bool {
	for _, t := range s.Topics {
		if topic == t {
			return true
		}
	}
	return false
}

type brokerActor struct {
}

func newBrokerActor() actor.Actor {
	return &brokerActor{}
}

func (state *brokerActor) Commands(typ string) interface{} {
	switch typ {
	case "command.subscribe":
		return &Subscribe{}
	case "command.unsubscribe":
		return &Unsubscribe{}
	}
	return nil
}

func (state *brokerActor) Events(msg interface{}) string {
	switch msg.(type) {
	case *messages.SubscribeSuccess:
		return "event.subscribe.success"
	case *messages.UnsubscribeSuccess:
		return "event.unsubscribe.success"
	}
	return ""
}

func (state *brokerActor) Receive(ctx actor.Context) {
}
