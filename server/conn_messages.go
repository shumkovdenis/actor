package server

import "github.com/shumkovdenis/club/server/core"

type ConnMessage interface {
	ConnMessage()
}

type ConnReadFailed struct{}

func (*ConnReadFailed) ConnMessage() {}

func (*ConnReadFailed) Event() string {
	return "event.conn.failed"
}

func (*ConnReadFailed) Code() string {
	return "conn_read_failed"
}

type Login struct {
	SessionID string `mapstructure:"session_id"`
}

func (*Login) ConnMessage() {}

func (*Login) Command() string {
	return "command.login"
}

type LoginSuccess struct{}

func (*LoginSuccess) ConnMessage() {}

func (*LoginSuccess) Event() string {
	return "event.login.success"
}

type LoginFailed struct {
	code core.Code
}

func (*LoginFailed) ConnMessage() {}

func (*LoginFailed) Event() string {
	return "event.login.failed"
}

func (m *LoginFailed) Code() string {
	if m.code != nil {
		return m.code.Code()
	}
	return ""
}

type AlreadyLogged struct{}

func (*AlreadyLogged) ConnMessage() {}

func (*AlreadyLogged) Event() string {
	return "event.already_logged"
}

type NotLogged struct{}

func (*NotLogged) ConnMessage() {}

func (*NotLogged) Event() string {
	return "event.not_logged"
}
