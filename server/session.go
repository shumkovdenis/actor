package server

// type SessionError struct {
// 	Code int
// }

// func newSessionError(code int) *SessionError {
// 	return &SessionError{code}
// }

// func (e *SessionError) Error() string {
// 	return fmt.Sprintf("session error: %d", e.Code)
// }

type Session struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
}
