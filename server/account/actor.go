package account

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/config"
)

type accountActor struct {
	sessionPID         *actor.PID
	username           string
	password           string
	tickerBalance      *time.Ticker
	tickerJackpotsTops *time.Ticker
	tickerJackpotsList *time.Ticker
}

func NewActor() actor.Actor {
	return &accountActor{}
}

func (state *accountActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.sessionPID = ctx.Parent()
	case *Authorize:
		res := authorize(msg.Username, msg.Password)
		ctx.Respond(res)
		if _, ok := res.(*Authorized); ok {
			state.username = msg.Username
			state.password = msg.Password
			ctx.SetBehavior(state.authorized)
		}
	case Message:
		ctx.Respond(&NotAuthorized{})
	}
}

func (state *accountActor) authorized(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Authorize:
		ctx.Respond(&AlreadyAuthorized{})
	case *GetBalance:
		res := getBalance(state.username, state.password)
		ctx.Respond(res)
	case *GetGameSession:
		res := getGameSession(state.username, state.password, msg.GameID)
		ctx.Respond(res)
	case *Withdraw:
		res := withdraw(state.username, state.password)
		ctx.Respond(res)
		if _, ok := res.(*WithdrawSuccess); ok {
			ctx.SetBehavior(state.Receive)
		}
	case *Cashback:
		res := cashback(state.username, state.password)
		ctx.Respond(res)
	case *GetJackpotsTops:
		res := jackpotsTops(state.username, state.password)
		ctx.Respond(res)
	case *GetJackpotsList:
		res := jackpotsList(state.username, state.password)
		ctx.Respond(res)
	case *StartLiveBalance:
		state.startLiveBalance()
	case *StopLiveBalance:
		state.stopLiveBalance()
	case *StartLiveJackpotsTops:
		state.startLiveJackpotsTops()
	case *StopLiveJackpotsTops:
		state.stopLiveJackpotsTops()
	case *StartLiveJackpotsList:
		state.startLiveJackpotsList()
	case *StopLiveJackpotsList:
		state.stopLiveJackpotsList()
	}
}

func (state *accountActor) startLiveBalance() {
	if state.tickerBalance != nil {
		return
	}

	conf := config.AccountAPI()

	state.tickerBalance = time.NewTicker(conf.BalanceInterval * time.Millisecond)

	go func() {
		for _ = range state.tickerBalance.C {
			res := getBalance(state.username, state.password)
			state.sessionPID.Tell(res)
		}
	}()
}

func (state *accountActor) stopLiveBalance() {
	if state.tickerBalance != nil {
		state.tickerBalance.Stop()
		state.tickerBalance = nil
	}
}

func (state *accountActor) startLiveJackpotsTops() {
	if state.tickerJackpotsTops != nil {
		return
	}

	conf := config.AccountAPI()

	state.tickerJackpotsTops = time.NewTicker(conf.JackpotsTopsInterval * time.Millisecond)

	go func() {
		for _ = range state.tickerJackpotsTops.C {
			res := jackpotsTops(state.username, state.password)
			state.sessionPID.Tell(res)
		}
	}()
}

func (state *accountActor) stopLiveJackpotsTops() {
	if state.tickerJackpotsTops != nil {
		state.tickerJackpotsTops.Stop()
		state.tickerJackpotsTops = nil
	}
}

func (state *accountActor) startLiveJackpotsList() {
	if state.tickerJackpotsList != nil {
		return
	}

	conf := config.AccountAPI()

	state.tickerJackpotsList = time.NewTicker(conf.JackpotsListInterval * time.Millisecond)

	go func() {
		for _ = range state.tickerJackpotsList.C {
			res := jackpotsList(state.username, state.password)
			state.sessionPID.Tell(res)
		}
	}()
}

func (state *accountActor) stopLiveJackpotsList() {
	if state.tickerJackpotsList != nil {
		state.tickerJackpotsList.Stop()
		state.tickerJackpotsList = nil
	}
}
