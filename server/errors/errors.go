package errors

var errID uint32

func nextErrID() uint32 {
	errID++
	return errID
}

type Err struct {
	ID    uint32
	Code  string
	Cause error
}

func (e *Err) Error() string {
	err := e.Previous
	switch {
	case err == nil:
		return ""
	}
	return ""
}

func NewErr(code string) Err {
	return Err{
		ID:   nextErrID(),
		Code: code,
	}
}

type ErrAPICreateSession struct {
	Err
}

func NewErrAPICreateSession() error {
	return &ErrAPICreateSession{NewErr("api_create_session_fail")}
}

type ErrRoomNotFound struct {
	Err
}

func NewErrRoomNotFound() error {
	return &ErrRoomNotFound{NewErr("room_not_found")}
}

func Wrap(previous, cause error) error {
	err := &Err{
		Previous: previous,
		Cause:    cause,
	}
	return err
}

func Mask(previous error) error {
	return &Err{
		Previous: previous,
	}
}
