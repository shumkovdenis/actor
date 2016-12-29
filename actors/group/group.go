package group

import (
	"github.com/AsynkronIT/gam/actor"
	"github.com/emirpasic/gods/sets/hashset"
)

type Use struct {
	Producer  *actor.PID
	Validator ValidateMessageFunc
}

type Join struct {
	Consumer *actor.PID
}

type Joined struct {
	Consumer *actor.PID
}

type Leave struct {
	Consumer *actor.PID
}

type Left struct {
	Consumer *actor.PID
}

type ValidateMessageFunc func(msg interface{}) bool

type groupActor struct {
	producer  *actor.PID
	consumers *hashset.Set
	validator ValidateMessageFunc
}

func NewActor() actor.Actor {
	return &groupActor{consumers: hashset.New()}
}

func (state *groupActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Use:
		state.producer = msg.Producer
		state.validator = msg.Validator
		ctx.Become(state.used)
	}
}

func (state *groupActor) used(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Join:
		state.consumers.Add(msg.Consumer)
		state.producer.Tell(&Joined{msg.Consumer})
	case *Leave:
		state.consumers.Remove(msg.Consumer)
		state.producer.Tell(&Left{msg.Consumer})
	default:
		if state.validator(msg) {
			for _, consumer := range state.consumers.Values() {
				consumer.(*actor.PID).Tell(msg)
			}
		}
	}
}
