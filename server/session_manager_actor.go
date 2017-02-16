package server

import (
	"go.uber.org/zap"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/maps/treemap"
)

const (
	CreateSessionFail = "create_session_fail"
	SessionNotFound   = "session_not_found"
)

type CreateSession struct {
	SessionID string
	RoomID    string `json:"room_id"`
}

type CreateSessionSuccess struct {
	Session *Session
}

type GetSession struct {
	SessionID string
}

type GetSessionSuccess struct {
	SessionPID *actor.PID
}

type sessionManagerActor struct {
	sessions       *treemap.Map
	roomManagerPID *actor.PID
}

func newSessionManagerActor() actor.Actor {
	return &sessionManagerActor{
		sessions: treemap.NewWithStringComparator(),
	}
}

func (state *sessionManagerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.roomManagerPID = actor.NewLocalPID("rooms")
	case *CreateSession:
		getRoom := &GetRoom{
			RoomID: msg.RoomID,
		}

		future := state.roomManagerPID.RequestFuture(getRoom, Timeout)
		res, err := future.Result()
		if err != nil {
			log.Error("create session fail: get room fail",
				zap.Error(err),
			)
			ctx.Respond(newFail(CreateSessionFail))
			return
		}

		if fail, ok := res.(Fail); ok {
			ctx.Respond(fail)
			return
		}

		roomPID := res.(*GetRoomSuccess).RoomPID

		props := actor.FromInstance(newSessionActor(roomPID))
		sessionPID, _ := ctx.SpawnNamed(props, msg.SessionID)

		joinRoom := &JoinRoom{
			SessionPID: sessionPID,
		}

		future = roomPID.RequestFuture(joinRoom, Timeout)
		res, err = future.Result()
		if err != nil {
			log.Error("create session fail: join room fail",
				zap.Error(err),
			)
			ctx.Respond(newFail(CreateSessionFail))
			return
		}

		if fail, ok := res.(Fail); ok {
			ctx.Respond(fail)
			return
		}

		state.sessions.Put(msg.SessionID, sessionPID)

		ctx.Respond(&CreateSessionSuccess{
			Session: &Session{
				ID:     msg.SessionID,
				RoomID: msg.RoomID,
			},
		})
	case *GetSession:
		pid, ok := state.sessions.Get(msg.SessionID)
		if !ok {
			ctx.Respond(newFail(SessionNotFound))
			return
		}

		ctx.Respond(&GetSessionSuccess{
			SessionPID: pid.(*actor.PID),
		})
	}
}
