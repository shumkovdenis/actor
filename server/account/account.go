package account

import (
	"github.com/shumkovdenis/club/logger"
)

var log = logger.Get()

type Incoming interface {
	Account()
}

type Outgoing interface {
	Account()
}

type Authorize struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (c *Authorize) Account() {}

func (*Authorize) Command() string {
	return "command.account.authorize"
}

type Authorized struct {
	Categories []Category `json:"categories"`
}

func (e *Authorized) Account() {}

func (*Authorized) Event() string {
	return "event.account.authorized"
}

type AlreadyAuthorized struct {
}

func (e *AlreadyAuthorized) Account() {}

func (e *AlreadyAuthorized) Event() string {
	return "event.account." + e.Code()
}

func (*AlreadyAuthorized) Code() string {
	return "already_authorized"
}

type NotAuthorized struct{}

func (e *NotAuthorized) Account() {}

func (e *NotAuthorized) Event() string {
	return "event.account." + e.Code()
}

func (*NotAuthorized) Code() string {
	return "not_authorized"
}

type AuthorizationFailed struct{}

func (e *AuthorizationFailed) Account() {}

func (e *AuthorizationFailed) Event() string {
	return "event.account." + e.Code()
}

func (*AuthorizationFailed) Code() string {
	return "authorization_failed"
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
