package server

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
)

type accountActor struct {
	account  string
	password string
}

func newAccountActor() actor.Actor {
	return &accountActor{}
}

func (state *accountActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *AccountAuth:
		proc := accountAuth(msg.Account, msg.Password)
		if Process(proc, ctx.Respond) {
			state.account = msg.Account
			state.password = msg.Password
			ctx.SetBehavior(state.authorized)
		}
	case *AccountBalance:
		err := newErr(ErrAccountNotAuth).LogErr()
		err = newErr(ErrAccountBalance).Wrap(err).LogErr()
		ctx.Respond(&AccountBalanceFail{err})
	case *AccountSession:
		err := newErr(ErrAccountNotAuth).LogErr()
		err = newErr(ErrAccountSession).Wrap(err).LogErr()
		ctx.Respond(&AccountSessionFail{err})
	case *AccountWithdraw:
		err := newErr(ErrAccountNotAuth).LogErr()
		err = newErr(ErrAccountWithdraw).Wrap(err).LogErr()
		ctx.Respond(&AccountWithdrawFail{err})
	}
}

func (state *accountActor) authorized(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *AccountAuth:
		err := newErr(ErrAccountAlreadyAuth).LogErr()
		err = newErr(ErrAccountAuth).Wrap(err).LogErr()
		ctx.Respond(&AccountAuthFail{err})
	case *AccountBalance:
		proc := accountBalance(state.account, state.password)
		Process(proc, ctx.Respond)
	case *AccountSession:
		proc := accountSession(state.account, state.password, msg.GameID)
		Process(proc, ctx.Respond)
	case *AccountWithdraw:
		proc := accountWithdraw(state.account, state.password)
		if Process(proc, ctx.Respond) {
			ctx.SetBehavior(state.Receive)
		}
	}
}

func accountAuth(account, password string) Proc {
	return func(tell Tell) bool {
		conf := config.AccountAPI()

		resp, err := resty.R().
			SetFormData(map[string]string{
				"auth_submit":   conf.Type + "_CLIENT_AUTH",
				"auth_username": account,
				"auth_password": password,
			}).
			Post(conf.URL)
		if err != nil {
			err := newErr(ErrAccountAuth).Error(err).LogErr()
			tell(&AccountAuthFail{err})
			return false
		}

		res := &struct {
			Result string `json:"result"`
			Code   int    `json:"code"`
			Groups []struct {
				Title string `json:"title"`
				Games []struct {
					ID    string `json:"id"`
					Title string `json:"title"`
				}
			} `json:"groups"`
		}{}
		if err = json.Unmarshal(resp.Body(), res); err != nil {
			err := newErr(ErrAccountAuth).Error(err).LogErr()
			tell(&AccountAuthFail{err})
			return false
		}

		if res.Result == "Error" {
			e := fmt.Errorf("error: %d", res.Code)
			err := newErr(ErrAccountAuth).Error(e).LogErr()
			tell(&AccountAuthFail{err})
			return false
		}

		categories := make([]Category, len(res.Groups))
		for i, group := range res.Groups {
			games := make([]Game, len(group.Games))
			for j, game := range group.Games {
				games[j] = Game{
					ID:    game.ID,
					Title: game.Title,
				}
			}
			categories[i] = Category{
				Title: group.Title,
				Games: games,
			}
		}

		tell(&AccountAuthSuccess{categories})
		return true
	}
}

func accountBalance(account, password string) Proc {
	return func(tell Tell) bool {
		conf := config.AccountAPI()

		resp, err := resty.R().
			SetFormData(map[string]string{
				"auth_submit":   conf.Type + "_CLIENT_BALANCE",
				"auth_username": account,
				"auth_password": password,
			}).
			Post(conf.URL)
		if err != nil {
			err := newErr(ErrAccountBalance).Error(err).LogErr()
			tell(&AccountBalanceFail{err})
			return false
		}

		res := &struct {
			Result  string  `json:"result"`
			Code    int     `json:"code"`
			Balance float64 `json:"balance"`
		}{}
		if err = json.Unmarshal(resp.Body(), res); err != nil {
			err := newErr(ErrAccountBalance).Error(err).LogErr()
			tell(&AccountBalanceFail{err})
			return false
		}

		if res.Result == "Error" {
			e := fmt.Errorf("error: %d", res.Code)
			err := newErr(ErrAccountBalance).Error(e).LogErr()
			tell(&AccountBalanceFail{err})
			return false
		}

		tell(&AccountBalanceSuccess{res.Balance})
		return true
	}
}

func accountSession(account, password string, gameID int) Proc {
	return func(tell Tell) bool {
		conf := config.AccountAPI()

		resp, err := resty.R().
			SetFormData(map[string]string{
				"auth_submit":   conf.Type + "_CLIENT_SESSION",
				"auth_username": account,
				"auth_password": password,
				"game_id":       strconv.Itoa(gameID),
			}).
			Post(conf.URL)
		if err != nil {
			err := newErr(ErrAccountSession).Error(err).LogErr()
			tell(&AccountSessionFail{err})
			return false
		}

		res := &struct {
			Result      string `json:"result"`
			Code        int    `json:"code"`
			SessionUUID string `json:"session_uuid"`
			GameUUID    string `json:"game_uuid"`
			Host        string `json:"host"`
		}{}
		if err = json.Unmarshal(resp.Body(), res); err != nil {
			err := newErr(ErrAccountSession).Error(err).LogErr()
			tell(&AccountSessionFail{err})
			return false
		}

		if res.Result == "Error" {
			e := fmt.Errorf("error: %d", res.Code)
			err := newErr(ErrAccountSession).Error(e).LogErr()
			tell(&AccountSessionFail{err})
			return false
		}

		tell(&AccountSessionSuccess{
			SessionID: res.SessionUUID,
			GameID:    res.GameUUID,
			ServerURL: res.Host,
		})
		return true
	}
}

func accountWithdraw(account, password string) Proc {
	return func(tell Tell) bool {
		conf := config.AccountAPI()

		resp, err := resty.R().
			SetFormData(map[string]string{
				"auth_submit":   conf.Type + "_CLIENT_WITHDRAW",
				"auth_username": account,
				"auth_password": password,
			}).
			Post(conf.URL)
		if err != nil {
			err := newErr(ErrAccountWithdraw).Error(err).LogErr()
			tell(&AccountWithdrawFail{err})
			return false
		}

		res := &struct {
			Result string `json:"result"`
			Code   int    `json:"code"`
		}{}
		if err = json.Unmarshal(resp.Body(), res); err != nil {
			err := newErr(ErrAccountWithdraw).Error(err).LogErr()
			tell(&AccountWithdrawFail{err})
			return false
		}

		if res.Result == "Error" {
			e := fmt.Errorf("error: %d", res.Code)
			err := newErr(ErrAccountWithdraw).Error(e).LogErr()
			tell(&AccountWithdrawFail{err})
			return false
		}

		tell(&AccountWithdrawSuccess{})
		return true
	}
}

type AccountAuth struct {
	Account  string `mapstructure:"account"`
	Password string `mapstructure:"password"`
}

func (*AccountAuth) Command() string {
	return "command.account.auth"
}

type AccountAuthSuccess struct {
	Categories []Category `json:"categories"`
}

func (*AccountAuthSuccess) Event() string {
	return "event.account.auth.success"
}

type AccountAuthFail struct {
	*Err
}

func (*AccountAuthFail) Event() string {
	return "event.account.auth.fail"
}

type AccountBalance struct {
}

func (*AccountBalance) Command() string {
	return "command.account.balance"
}

type AccountBalanceSuccess struct {
	Balance float64 `json:"balance"`
}

func (*AccountBalanceSuccess) Event() string {
	return "event.account.balance.success"
}

type AccountBalanceFail struct {
	*Err
}

func (*AccountBalanceFail) Event() string {
	return "event.account.balance.fail"
}

type AccountSession struct {
	GameID int `mapstructure:"game_id"`
}

func (*AccountSession) Command() string {
	return "command.account.session"
}

type AccountSessionSuccess struct {
	SessionID string `json:"session_id"`
	GameID    string `json:"game_id"`
	ServerURL string `json:"server_url"`
}

func (*AccountSessionSuccess) Event() string {
	return "event.account.session.success"
}

type AccountSessionFail struct {
	*Err
}

func (*AccountSessionFail) Event() string {
	return "event.account.session.fail"
}

type AccountWithdraw struct {
}

func (*AccountWithdraw) Command() string {
	return "command.account.withdraw"
}

type AccountWithdrawSuccess struct {
}

func (*AccountWithdrawSuccess) Event() string {
	return "event.account.withdraw.success"
}

type AccountWithdrawFail struct {
	*Err
}

func (*AccountWithdrawFail) Event() string {
	return "event.account.withdraw.fail"
}

type Category struct {
	Title string `json:"title"`
	Games []Game `json:"games"`
}

type Game struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Thumb string `json:"thumb"`
}
