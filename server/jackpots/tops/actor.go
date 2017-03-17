package tops

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/shumkovdenis/club/config"
)

type topsActor struct {
	sessions *hashset.Set
	ticker   *time.Ticker
}

func NewActor() actor.Actor {
	return &topsActor{
		sessions: hashset.New(),
	}
}

func (state *topsActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Join:
		state.sessions.Add(msg.SessionPID)
		if state.sessions.Size() == 1 {
			state.start()
		}
	case *Leave:
		state.sessions.Remove(msg.SessionPID)
		if state.sessions.Size() == 0 {
			state.stop()
		}
	case *Get:
		res := tops()
		ctx.Respond(res)
	}
}

func (state *topsActor) start() {
	conf := config.AccountAPI()

	state.ticker = time.NewTicker(conf.JackpotsTopsInterval * time.Millisecond)

	go func() {
		for _ = range state.ticker.C {
			res := tops()
			state.tell(res)
		}
	}()
}

func (state *topsActor) stop() {
	if state.ticker != nil {
		state.ticker.Stop()
		state.ticker = nil
	}
}

func (state *topsActor) tell(msg interface{}) {
	for _, pid := range state.sessions.Values() {
		pid.(*actor.PID).Tell(msg)
	}
}

func tops() Message {
	tops := make([]Jackpot, 5)
	tops[0] = Jackpot{"1191100006", 12468.00, 1489754067}
	tops[1] = Jackpot{"1191100006", 12468.00, 1489754067}
	tops[2] = Jackpot{"1191100006", 12468.00, 1489754067}
	tops[3] = Jackpot{"1191100006", 12468.00, 1489754067}
	tops[4] = Jackpot{"1191100006", 12468.00, 1489754067}
	return &Tops{tops}
}
