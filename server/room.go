package server

// type RoomError struct {
// 	Code int
// }

// func newRoomError(code int) *RoomError {
// 	return &RoomError{code}
// }

// func (e *RoomError) Error() string {
// 	return fmt.Sprintf("room error: %d", e.Code)
// }

type Room struct {
	ID string `json:"id"`
	// Sessions []string `json:"sessions"`
	// maxSessions int
	// sessions    *hashset.Set
}

// func newRoom() *Room {
// 	id := uuid.NewV4().String()
// 	return &Room{
// 		ID:          id,
// 		maxSessions: 2,
// 		sessions:    hashset.New(),
// 	}
// }

// func (r *Room) Join(s *Session) error {
// 	if r.Full() {
// 		log.Warn("join room: room full",
// 			zap.String("room", r.ID),
// 		)

// 		return &RoomError{RoomFull}
// 	}

// 	if r.sessions.Contains(s) {
// 		log.Warn("join room: room contains session",
// 			zap.String("room", r.ID),
// 			zap.String("session", s.ID),
// 		)

// 		return &RoomError{RoomContainsSession}
// 	}

// 	r.sessions.Add(s)

// 	return nil
// }

// func (r *Room) Used() bool {
// 	return r.sessions.Size() > 0
// }

// func (r *Room) Full() bool {
// 	return r.maxSessions == r.sessions.Size()
// }
