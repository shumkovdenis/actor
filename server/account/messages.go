package account

import (
	"github.com/shumkovdenis/club/logger"
)

var log = logger.Get()

type Message interface {
	AccountMessage()
}

type Authorize struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (*Authorize) AccountMessage() {}

func (*Authorize) Command() string {
	return "command.account.authorize"
}

type Authorized struct {
	Categories []Category `json:"categories"`
}

func (*Authorized) AccountMessage() {}

func (*Authorized) Event() string {
	return "event.account.authorized"
}

type AlreadyAuthorized struct{}

func (*AlreadyAuthorized) AccountMessage() {}

func (m *AlreadyAuthorized) Event() string {
	return "event.account.already_authorized"
}

type NotAuthorized struct{}

func (*NotAuthorized) AccountMessage() {}

func (m *NotAuthorized) Event() string {
	return "event.account.not_authorized"
}

type AuthorizationFailed struct{}

func (*AuthorizationFailed) AccountMessage() {}

func (m *AuthorizationFailed) Event() string {
	return "event.account.authorize.failed"
}

type GetBalance struct{}

func (*GetBalance) AccountMessage() {}

func (*GetBalance) Command() string {
	return "command.account.balance"
}

type Balance struct {
	Balance float64 `json:"balance"`
}

func (*Balance) AccountMessage() {}

func (*Balance) Event() string {
	return "event.account.balance"
}

type GetBalanceFailed struct{}

func (*GetBalanceFailed) AccountMessage() {}

func (*GetBalanceFailed) Event() string {
	return "event.account.balance.failed"
}

type GetGameSession struct {
	GameID int `mapstructure:"game_id"`
}

func (*GetGameSession) AccountMessage() {}

func (*GetGameSession) Command() string {
	return "command.account.session"
}

type GameSession struct {
	SessionID string `json:"session_id"`
	GameID    string `json:"game_id"`
	ServerURL string `json:"server_url"`
}

func (*GameSession) AccountMessage() {}

func (*GameSession) Event() string {
	return "event.account.session"
}

type GetGameSessionFailed struct{}

func (*GetGameSessionFailed) AccountMessage() {}

func (*GetGameSessionFailed) Event() string {
	return "event.account.session.failed"
}

type Withdraw struct{}

func (*Withdraw) AccountMessage() {}

func (*Withdraw) Command() string {
	return "command.account.withdraw"
}

type WithdrawSuccess struct{}

func (*WithdrawSuccess) AccountMessage() {}

func (*WithdrawSuccess) Event() string {
	return "event.account.withdraw"
}

type WithdrawFailed struct{}

func (*WithdrawFailed) AccountMessage() {}

func (*WithdrawFailed) Event() string {
	return "event.account.withdraw.failed"
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
