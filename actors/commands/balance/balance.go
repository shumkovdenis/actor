package balance

import (
	"time"

	"github.com/AsynkronIT/gam/actor"
)

type Balance struct {
}

type BalanceSuccess struct {
}

type BalanceFail struct {
}

type balanceActor struct {
}

func NewActor() actor.Actor {
	return &balanceActor{}
}

func (state *balanceActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		go func() {
			tick := time.Tick(5 * time.Second)
			for _ = range tick {
				ctx.Parent().Tell(&BalanceSuccess{})
			}
		}()
	case *Balance:
		if err := fetch(msg); err != nil {
		}
		ctx.Respond(&BalanceSuccess{})
	}
}

func fetch(auth *Balance) error {
	return nil
}
