package server

import "github.com/shumkovdenis/club/server/core"

type BrokerMessage interface {
	BrokerMessage()
}

type Subscribe struct {
	Topics []string `mapstructure:"topics"`
}

func (*Subscribe) BrokerMessage() {}

func (*Subscribe) Command() string {
	return "command.subscribe"
}

func (m *Subscribe) Contains(evt core.Event) bool {
	for _, topic := range m.Topics {
		if topic == evt.Event() {
			return true
		}
	}
	return false
}

type SubscribeSuccess struct {
	Topics []string `json:"topics"`
}

func (*SubscribeSuccess) BrokerMessage() {}

func (*SubscribeSuccess) Event() string {
	return "event.subscribe.success"
}

type SubscribeFailed struct {
	Message string `json:"message"`
}

func (*SubscribeFailed) BrokerMessage() {}

func (*SubscribeFailed) Event() string {
	return "event.subscribe.failed"
}

type Unsubscribe struct {
	Topics []string `mapstructure:"topics"`
}

func (*Unsubscribe) BrokerMessage() {}

func (*Unsubscribe) Command() string {
	return "command.unsubscribe"
}

func (m *Unsubscribe) Contains(evt core.Event) bool {
	for _, topic := range m.Topics {
		if topic == evt.Event() {
			return true
		}
	}
	return false
}

type UnsubscribeSuccess struct {
	Topics []string `json:"topics"`
}

func (*UnsubscribeSuccess) BrokerMessage() {}

func (*UnsubscribeSuccess) Event() string {
	return "event.unsubscribe.success"
}

type UnsubscribeFailed struct {
	Message string `json:"message"`
}

func (*UnsubscribeFailed) BrokerMessage() {}

func (*UnsubscribeFailed) Event() string {
	return "event.unsubscribe.failed"
}
