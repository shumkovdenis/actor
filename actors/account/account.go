package account

import (
	"github.com/AsynkronIT/gam/actor"
	"github.com/shumkovdenis/club/manifest"
)

// Auth -> command.account.auth
type Auth struct {
	Account  string `mapstructure:"account"`
	Password string `mapstructure:"password"`
}

// AuthSuccess -> event.account.auth.success
type AuthSuccess struct {
	Categories []Category `json:"categories"`
}

// AuthFail -> event.account.auth.fail
type AuthFail struct {
	Message string `json:"message"`
}

// Balance -> command.account.balance
type Balance struct {
}

// BalanceSuccess -> event.account.balance.success
type BalanceSuccess struct {
	Balance float64 `json:"balance"`
}

// BalanceFail -> event.account.balance.fail
type BalanceFail struct {
	Message string `json:"message"`
}

// Session -> command.account.session
type Session struct {
	GameID int `mapstructure:"game_id"`
}

// SessionSuccess -> event.account.session.success
type SessionSuccess struct {
	SessionID string `json:"session_id"`
	GameID    string `json:"game_id"`
	ServerURL string `json:"server_url"`
}

// SessionFail -> event.account.session.fail
type SessionFail struct {
	Message string `json:"message"`
}

// Withdraw -> command.account.withdraw
type Withdraw struct {
}

// WithdrawSuccess -> event.account.withdraw.success
type WithdrawSuccess struct {
}

// WithdrawFail -> event.account.withdraw.fail
type WithdrawFail struct {
	Message string `json:"message"`
}

// Fail -> event.account.fail
type Fail struct {
	Message string `json:"message"`
}

type accountActor struct {
	urlAPI   string
	typeAPI  string
	account  string
	password string
}

func NewActor() actor.Actor {
	conf := manifest.Get().Config.AccountAPI
	return &accountActor{
		urlAPI:  conf.URL,
		typeAPI: conf.Type,
	}
}

func (state *accountActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		ctx.Become(state.started)
	}
}

func (state *accountActor) started(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Auth:
		success, err := state.auth(msg.Account, msg.Password)
		if err != nil {
			ctx.Respond(&AuthFail{err.Error()})
			return
		}
		state.account = msg.Account
		state.password = msg.Password
		ctx.Respond(success)
		ctx.Become(state.authorized)
	case *Balance, *Session, *Withdraw:
		ctx.Respond(&Fail{"Account is not authorized"})
	}
}

func (state *accountActor) authorized(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Balance:
		success, err := state.balance()
		if err != nil {
			ctx.Respond(&BalanceFail{err.Error()})
			return
		}
		ctx.Respond(success)
	case *Session:
		success, err := state.session(msg.GameID)
		if err != nil {
			ctx.Respond(&SessionFail{err.Error()})
			return
		}
		ctx.Respond(success)
	case *Withdraw:
		success, err := state.withdraw()
		if err != nil {
			ctx.Respond(&WithdrawFail{err.Error()})
			return
		}
		ctx.Respond(success)
		ctx.Become(state.started)
	case *Auth:
		ctx.Respond(&Fail{"Account already authorized"})
	}
}
