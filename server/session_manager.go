package server

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
	uuid "github.com/satori/go.uuid"
	"github.com/uber-go/zap"
)

const (
	SessionNotFound = iota
	SessionUsed
)

type SessionError struct {
	Code int
}

func (e *SessionError) Error() string {
	return fmt.Sprintf("session manager error: %d", e.Code)
}

type SessionConf struct {
	RoomID string `json:"room_id"`
}

type Session struct {
	ID   string       `json:"id"`
	Used bool         `json:"used"`
	Conf *SessionConf `json:"conf"`
}

type SessionManager interface {
	CreateSession(conf *SessionConf) (*Session, error)
	UseSession(id string) error
}

type sessionManager struct {
	sessions *treemap.Map
}

func newSessionManager() SessionManager {
	return &sessionManager{
		sessions: treemap.NewWithStringComparator(),
	}
}

func (m *sessionManager) CreateSession(conf *SessionConf) (*Session, error) {
	id := uuid.NewV4().String()

	session := &Session{
		ID:   id,
		Conf: conf,
	}

	m.sessions.Put(id, session)

	return session, nil
}

func (m *sessionManager) UseSession(id string) error {
	value, ok := m.sessions.Get(id)
	if !ok {
		log.Warn("use session: session not found", zap.String("session", id))

		return &SessionError{SessionNotFound}
	}

	session := value.(Session)

	if session.Used {
		log.Warn("use session: session used", zap.String("session", id))

		return &SessionError{SessionUsed}
	}

	session.Used = true

	return nil
}
