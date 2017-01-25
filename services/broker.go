package services

import (
	"github.com/emirpasic/gods/sets/treeset"
)

type Broker interface {
	Topics() *treeset.Set
	AddTopics(topics []string)
	RemoveTopics(topics []string)
	Contains(topic string) bool
}

type broker struct {
	topics *treeset.Set
}

func NewBroker() Broker {
	return &broker{
		topics: treeset.NewWithStringComparator(),
	}
}

func (b *broker) Topics() *treeset.Set {
	return b.topics
}

func (b *broker) AddTopics(topics []string) {
	for _, topic := range topics {
		b.topics.Add(topic)
	}
}

func (b *broker) RemoveTopics(topics []string) {
	for _, topic := range topics {
		b.topics.Remove(topic)
	}
}

func (b *broker) Contains(topic string) bool {
	return b.topics.Contains(topic)
}
