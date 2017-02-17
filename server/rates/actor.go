package rates

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
)

type ratesActor struct {
	sessions *hashset.Set
	ticker   *time.Ticker
}

func NewActor() actor.Actor {
	return &ratesActor{
		sessions: hashset.New(),
	}
}

func (state *ratesActor) Receive(ctx actor.Context) {
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
	}
}

func (state *ratesActor) start() {
	conf := config.RatesAPI()

	state.ticker = time.NewTicker(conf.GetInterval * time.Millisecond)

	go func() {
		for _ = range state.ticker.C {
			res := rates()
			state.tell(res)
		}
	}()
}

func (state *ratesActor) stop() {
	state.ticker.Stop()
}

func (state *ratesActor) tell(msg interface{}) {
	for _, pid := range state.sessions.Values() {
		pid.(*actor.PID).Tell(msg)
	}
}

func rates() Message {
	conf := config.RatesAPI()

	resp, err := resty.R().Get(conf.URL)
	if err != nil {
		log.Error("rates failed",
			zap.Error(err),
		)
		return &RatesFailed{}
	}

	if resp.StatusCode() != http.StatusOK {
		log.Error("rates failed",
			zap.Int("code", resp.StatusCode()),
		)
		return &RatesFailed{}
	}

	res := make([]*Rate, 0)
	if err = json.Unmarshal(resp.Body(), &res); err != nil {
		log.Error("rates failed",
			zap.Error(err),
		)
		return &RatesFailed{}
	}

	return &Rates{res}
}
