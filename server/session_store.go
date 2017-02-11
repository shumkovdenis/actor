package server

import (
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/juju/errors"
)

type SessionStore interface {
	Add(session *Session) error
	GetByID(id string) (*Session, error)
}

type sessionStore struct {
	sessions *treemap.Map
}

func newSessionStore() SessionStore {
	return &sessionStore{
		sessions: treemap.NewWithStringComparator(),
	}
}

func (s *sessionStore) Add(session *Session) error {
	s.sessions.Put(session.ID, session)
	return nil
}

func (s *sessionStore) GetByID(id string) (*Session, error) {
	session, ok := s.sessions.Get(id)
	if !ok {
		return nil, newSessionNotFound()
	}
	return session.(*Session), nil
}

type sessionNotFound struct {
	errors.Err
}

func newSessionNotFound() error {
	err := &sessionNotFound{errors.NewErr("session not found")}
	err.SetLocation(1)
	return err
}
