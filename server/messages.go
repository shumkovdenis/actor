package server

type Command interface {
	Command() string
}

type Event interface {
	Event() string
}

type command struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Subscribe struct {
	Topics []string `mapstructure:"topics"`
}

func (*Subscribe) Command() string {
	return "command.subscribe"
}

type SubscribeSuccess struct {
	Topics []string `json:"topics"`
}

func (*SubscribeSuccess) Event() string {
	return "event.subscribe.success"
}

type SubscribeFail struct {
	Message string `json:"message"`
}

func (*SubscribeFail) Event() string {
	return "event.subscribe.fail"
}

type Unsubscribe struct {
	Topics []string `mapstructure:"topics"`
}

func (*Unsubscribe) Command() string {
	return "command.unsubscribe"
}

type UnsubscribeSuccess struct {
	Topics []string `json:"topics"`
}

func (*UnsubscribeSuccess) Event() string {
	return "event.unsubscribe.success"
}

type UnsubscribeFail struct {
	Message string `json:"message"`
}

func (*UnsubscribeFail) Event() string {
	return "event.unsubscribe.fail"
}

type Login struct {
	SessionID string `mapstructure:"session_id"`
}

func (*Login) Command() string {
	return "command.login"
}

type LoginSuccess struct {
}

func (*LoginSuccess) Event() string {
	return "event.login.success"
}

type LoginFail struct {
	// *Err
}

func (*LoginFail) Event() string {
	return "event.login.fail"
}

type LeaveRoom struct {
	RoomID string
}

type LeaveRoomSuccess struct{}

type LeaveRoomFail error

type Success struct{}

// type Fail struct {
// 	*Err
// }

// func (*Fail) Event() string {
// 	return "event.fail"
// }

type Fail interface {
	Code() string
}

type fail struct {
	CodeFail string `json:"code"`
}

func newFail(code string) Fail {
	return &fail{code}
}

func (f *fail) Code() string {
	return f.CodeFail
}
