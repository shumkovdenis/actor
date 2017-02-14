package server

import (
	"encoding/json"
	"fmt"

	"github.com/uber-go/zap"
)

const (
	_ uint16 = iota

	ErrReadJSON
	ErrWriteJSON

	ErrToMessage
	ErrFromMessage

	ErrAPICreateRoom
	ErrAPICreateSession

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

	ErrAccountNotAuth
	ErrAccountAlreadyAuth
	ErrAccountAuth
	ErrAccountBalance
	ErrAccountSession
	ErrAccountWithdraw
)

var errID uint32

func nextErrID() uint32 {
	errID++
	return errID
}

type Err struct {
	ID      uint32
	Code    uint16
	Details string
	Child   *Err
}

func newErr(code uint16) *Err {
	return &Err{
		ID:   nextErrID(),
		Code: code,
	}
}

func (e *Err) Wrap(err *Err) *Err {
	e.Child = err
	return e
}

func (e *Err) Error(err error) *Err {
	e.Details = err.Error()
	return e
}

func (e *Err) LogErr() *Err {
	fields := make([]zap.Field, 0, 4)
	fields = append(fields,
		zap.Uint64("id", uint64(e.ID)),
		zap.Uint64("code", uint64(e.Code)),
	)
	if e.Child != nil {
		fields = append(fields,
			zap.Uint64("child", uint64(e.Child.ID)),
		)
	}
	if len(e.Details) > 0 {
		fields = append(fields,
			zap.String("details", e.Details),
		)
	}
	log.Error(e.Message(), fields...)
	return e
}

func (e *Err) Message() string {
	var msg string
	switch e.Code {
	case ErrReadJSON:
		msg = "read json fail"
	case ErrWriteJSON:
		msg = "write json fail"
	case ErrToMessage:
		msg = "conv to message fail"
	case ErrFromMessage:
		msg = "conv from message fail"
	case ErrAPICreateRoom:
		msg = "api create room fail"
	case ErrAPICreateSession:
		msg = "api create session fail"
	case ErrCreateRoom:
		msg = "create room fail"
	case ErrGetRoom:
		msg = "get room fail"
	case ErrJoinRoom:
		msg = "join room fail"
	case ErrRoomNotFound:
		msg = "room not found"
	case ErrRoomFull:
		msg = "room full"
	case ErrCreateSession:
		msg = "create session fail"
	case ErrGetSession:
		msg = "get session fail"
	case ErrUseSession:
		msg = "use session fail"
	case ErrSessionNotFound:
		msg = "session not found"
	case ErrLogin:
		msg = "login fail"
	case ErrAccountNotAuth:
		msg = "account not auth"
	case ErrAccountAlreadyAuth:
		msg = "account already auth"
	case ErrAccountAuth:
		msg = "account auth fail"
	case ErrAccountBalance:
		msg = "account balance fail"
	case ErrAccountSession:
		msg = "account session fail"
	case ErrAccountWithdraw:
		msg = "account withdraw fail"
	default:
		msg = fmt.Sprintf("error %d", e.Code)
	}
	return msg
}

func (e *Err) MarshalJSON() ([]byte, error) {
	err := e
	for err.Child != nil {
		err = err.Child
	}
	return json.Marshal(&struct {
		ID      uint32 `json:"id"`
		Code    uint16 `json:"code"`
		Message string `json:"message"`
	}{
		ID:      err.ID,
		Code:    err.Code,
		Message: err.Message(),
	})
}
