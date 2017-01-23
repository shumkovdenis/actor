package group

import (
	"reflect"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/sets/hashset"
)

type Use struct {
	Producer *actor.PID
	Types    []interface{}
}

type Used struct {
	Producer *actor.PID
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

type groupActor struct {
	producers *hashset.Set
	consumers *hashset.Set
	types     *hashset.Set
}

func New() actor.Actor {
	return &groupActor{
		producers: hashset.New(),
		consumers: hashset.New(),
		types:     hashset.New(),
	}
}

func (state *groupActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Use:
		state.producers.Add(msg.Producer)
		for _, t := range msg.Types {
			state.types.Add(t)
		}
		ctx.Respond(&Used{msg.Producer})
	case *Join:
		state.consumers.Add(msg.Consumer)
		state.tellProducers(&Joined{msg.Consumer})
	case *Leave:
		state.consumers.Remove(msg.Consumer)
		state.tellProducers(&Left{msg.Consumer})
	default:
		state.tellConsumers(msg)
	}
}

func (state *groupActor) tellProducers(msg interface{}) {
	for _, producer := range state.producers.Values() {
		producer.(*actor.PID).Tell(msg)
	}
}

func (state *groupActor) tellConsumers(msg interface{}) {
	for _, t := range state.types.Values() {
		if reflect.TypeOf(msg).AssignableTo(reflect.TypeOf(t)) {
			for _, consumer := range state.consumers.Values() {
				consumer.(*actor.PID).Tell(msg)
			}
			break
		}
	}
}
