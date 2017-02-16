package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
)

type ratesActor struct {
	sessions *hashset.Set
	ticker   *time.Ticker
}

func newRatesActor() actor.Actor {
	return &ratesActor{
		sessions: hashset.New(),
	}
}

func (state *ratesActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *JoinRates:
		state.sessions.Add(msg.SessionPID)
		if state.sessions.Size() == 1 {
			state.start()
		}
	case *LeaveRates:
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
			rates(state.tell)
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

func rates(tell Tell) bool {
	conf := config.RatesAPI()

	resp, err := resty.R().Get(conf.URL)
	if err != nil {
		// err := newErr(ErrRates).Error(err).LogErr()
		// tell(&RatesFail{err})
		return false
	}

	if resp.StatusCode() != http.StatusOK {
		// e := fmt.Errorf("status code: %d", resp.StatusCode())
		// err := newErr(ErrRates).Error(e).LogErr()
		// tell(&RatesFail{err})
		return false
	}

	res := make([]*Rate, 0)
	if err = json.Unmarshal(resp.Body(), &res); err != nil {
		// 	err := newErr(ErrRates).Error(err).LogErr()
		// 	tell(&RatesFail{err})
		return false
	}

	// tell(&RatesChange{res})
	return true
}

type JoinRates struct {
	SessionPID *actor.PID
}

type LeaveRates struct {
	SessionPID *actor.PID
}

type RatesChange struct {
	Rates []*Rate `json:"rates"`
}

func (*RatesChange) Event() string {
	return "event.rates.change"
}

type RatesFail struct {
	// *Err
}

func (*RatesFail) Event() string {
	return "event.rates.fail"
}

type Rate struct {
	Timestamp uint64  `json:"timestamp"`
	Value     float64 `json:"value"`
}
