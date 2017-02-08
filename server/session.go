package server

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

const (
	SessionNotFound = iota
	SessionUsed
)

type SessionError struct {
	Code int
}

func newSessionError(code int) *SessionError {
	return &SessionError{code}
}

func (e *SessionError) Error() string {
	return fmt.Sprintf("session error: %d", e.Code)
}

type SessionConf struct {
	RoomID string `json:"room_id"`
}

func newSessionConf() *SessionConf {
	return &SessionConf{}
}

type Session struct {
	ID   string       `json:"id"`
	Conf *SessionConf `json:"conf"`
	used bool
}

func newSession(conf *SessionConf) *Session {
	id := uuid.NewV4().String()
	return &Session{
		ID:   id,
		Conf: conf,
	}
}
