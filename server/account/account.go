package account

import (
	"github.com/shumkovdenis/club/logger"
)

var log = logger.Get()

type Message interface {
	Account()
}

type Authorize struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (*Authorize) Account() {}

func (*Authorize) Command() string {
	return "command.account.authorize"
}

type Authorized struct {
	Categories []Category `json:"categories"`
}

func (*Authorized) Account() {}

func (*Authorized) Event() string {
	return "event.account.authorized"
}

type AlreadyAuthorized struct{}

func (*AlreadyAuthorized) Account() {}

func (m *AlreadyAuthorized) Event() string {
	return "event.account.authorize." + m.Code()
}

func (*AlreadyAuthorized) Code() string {
	return "already_authorized"
}

type NotAuthorized struct{}

func (*NotAuthorized) Account() {}

func (m *NotAuthorized) Event() string {
	return "event.account.authorize." + m.Code()
}

func (*NotAuthorized) Code() string {
	return "not_authorized"
}

type AuthorizationFailed struct{}

func (*AuthorizationFailed) Account() {}

func (m *AuthorizationFailed) Event() string {
	return "event.account.authorize." + m.Code()
}

func (*AuthorizationFailed) Code() string {
	return "authorization_failed"
}

type GetBalance struct{}

func (*GetBalance) Account() {}

func (*GetBalance) Command() string {
	return "command.account.balance"
}

type Balance struct {
	Balance float64 `json:"balance"`
}

func (*Balance) Account() {}

func (*Balance) Event() string {
	return "event.account.balance"
}

type GetBalanceFailed struct{}

func (*GetBalanceFailed) Account() {}

func (*GetBalanceFailed) Event() string {
	return "event.account.balance.failed"
}

func (*GetBalanceFailed) Code() string {
	return "get_balance_failed"
}

type GetGameSession struct {
	GameID int `mapstructure:"game_id"`
}

func (*GetGameSession) Account() {}

func (*GetGameSession) Command() string {
	return "command.account.session"
}

type GameSession struct {
	SessionID string `json:"session_id"`
	GameID    string `json:"game_id"`
	ServerURL string `json:"server_url"`
}

func (*GameSession) Account() {}

func (*GameSession) Event() string {
	return "event.account.session"
}

type GetGameSessionFailed struct{}

func (*GetGameSessionFailed) Account() {}

func (*GetGameSessionFailed) Event() string {
	return "event.account.session.failed"
}

func (*GetGameSessionFailed) Code() string {
	return "get_game_session_failed"
}

type Withdraw struct{}

func (*Withdraw) Account() {}

func (*Withdraw) Command() string {
	return "command.account.withdraw"
}

type WithdrawSuccess struct{}

func (*WithdrawSuccess) Account() {}

func (*WithdrawSuccess) Event() string {
	return "event.account.withdraw"
}

type WithdrawFailed struct{}

func (*WithdrawFailed) Account() {}

func (*WithdrawFailed) Event() string {
	return "event.account.withdraw.failed"
}

func (*WithdrawFailed) Code() string {
	return "withdraw_failed"
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
