package rates

import (
	"log"

	"github.com/AsynkronIT/gam/actor"

	"github.com/emirpasic/gods/sets/hashset"
)

type Sub struct {
}

// Change -> event.rates.change
type Change struct {
}

// Fail -> event.rates.fail
type Fail struct {
}

type ratesActor struct {
	subs *hashset.Set
}

func NewActor() actor.Actor {
	return &ratesActor{hashset.New()}
}

func (state *ratesActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *Sub:
		state.subs.Add(ctx.Sender())
		log.Printf("[rates actor] count = %d", state.subs.Size())
	}
}
