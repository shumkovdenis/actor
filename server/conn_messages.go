package server

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

type LoginFailed struct{}

func (*LoginFailed) ConnMessage() {}

func (*LoginFailed) Event() string {
	return "event.login.failed"
}
