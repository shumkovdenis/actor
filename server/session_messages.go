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

func (*CreateSessionFailed) Fail() string {
	return "create_session_failed"
}

type GetSession struct {
	SessionID string
}

func (*GetSession) SessionMessage() {}

type SessionNotFound struct{}

func (*SessionNotFound) SessionMessage() {}

func (*SessionNotFound) Fail() string {
	return "session_not_found"
}

type UseSession struct {
	ConnPID *actor.PID
}

func (*UseSession) SessionMessage() {}

type UseSessionSuccess struct{}

func (*UseSessionSuccess) SessionMessage() {}

type SessionAlreadyUse struct{}

func (*SessionAlreadyUse) SessionMessage() {}

func (*SessionAlreadyUse) Fail() string {
	return "session_already_use"
}
