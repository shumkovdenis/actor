package server

import "github.com/AsynkronIT/protoactor-go/actor"

type SessionMessage interface {
	SessionMessage()
}

type CreateSession struct {
	SessionID string
	RoomID    string
}

func (*CreateSession) SessionMessage() {}

type Session struct {
	ID string
}

func (*Session) SessionMessage() {}

type CreateSessionFailed struct{}

func (*CreateSessionFailed) SessionMessage() {}

func (*CreateSessionFailed) Code() string {
	return "create_session_failed"
}

type GetSession struct {
	SessionID string
}

func (*GetSession) SessionMessage() {}

type SessionNotFound struct{}

func (*SessionNotFound) SessionMessage() {}

func (*SessionNotFound) Code() string {
	return "session_not_found"
}

type UseSession struct {
	ConnPID *actor.PID
}

func (*UseSession) SessionMessage() {}

type UseSessionSuccess struct{}

func (*UseSessionSuccess) SessionMessage() {}

type SessionAlreadyUsed struct{}

func (*SessionAlreadyUsed) SessionMessage() {}

func (*SessionAlreadyUsed) Code() string {
	return "session_already_used"
}

type SessionNotUsed struct{}

func (*SessionNotUsed) SessionMessage() {}

func (*SessionNotUsed) Code() string {
	return "session_not_used"
}

type FreeSession struct{}

func (*FreeSession) SessionMessage() {}

type FreeSessionSuccess struct{}

func (*FreeSessionSuccess) SessionMessage() {}
