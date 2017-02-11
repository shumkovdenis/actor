package server

/*type sessionManager struct {
	sessions *treemap.Map
}

func newSessionManager() *sessionManager {
	return &sessionManager{
		sessions: treemap.NewWithStringComparator(),
	}
}

func (m *sessionManager) Create(conf *SessionConf) *Session {
	session := newSession(conf)

	m.sessions.Put(id, session)

	return session, nil
}

func (m *sessionManager) Get(id string) (*Session, error) {
	session, ok := m.rooms.Get(id)
	if !ok {
		log.Warn("get session: session not found",
			zap.String("session", id),
		)

		return nil, &SessionError{SessionNotFound}
	}

	return session.(*Session), nil
}
*/

/*
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
*/
