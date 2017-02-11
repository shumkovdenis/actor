package server

import (
	"fmt"
)

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
	*Error
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

type Fail struct {
	*Error
}

func (*Fail) Event() string {
	return "event.fail"
}

const (
	ErrCode = iota

	ErrReadJSON
	ErrWriteJSON

	ErrToMessage
	ErrFromMessage

	ErrCreateRoom
	ErrGetRoom
	ErrJoinRoom
	ErrRoomNotFound
	ErrRoomFull

	ErrCreateSession
	ErrGetSession
	ErrUseSession
	ErrSessionNotFound

	ErrLogin
)

type Error struct {
	Code    int
	Message string `json:"message"`
}

func newError(code int) *Error {
	return &Error{
		Code: code,
	}
}

func newErrorWrap(code int, err error) *Error {
	return &Error{
		Code:    code,
		Message: err.Error(),
	}
}

func (e *Error) Error() string {
	var s string
	switch e.Code {
	case ErrReadJSON:
		s = "read json fail"
	case ErrWriteJSON:
		s = "write json fail"
	case ErrToMessage:
		s = "conv to message fail"
	case ErrFromMessage:
		s = "conv from message fail"
	case ErrCreateRoom:
		s = "create room fail"
	case ErrGetRoom:
		s = "get room fail"
	case ErrJoinRoom:
		s = "join room fail"
	case ErrRoomNotFound:
		s = "room not found"
	case ErrRoomFull:
		s = "room full"
	case ErrCreateSession:
		s = "create session fail"
	case ErrGetSession:
		s = "get session fail"
	case ErrUseSession:
		s = "use session fail"
	case ErrSessionNotFound:
		s = "session not found"
	case ErrLogin:
		s = "login fail"
	default:
		s = fmt.Sprintf("error %d", e.Code)
	}

	if len(e.Message) > 0 {
		s += ": " + e.Message
	}

	return s
}
