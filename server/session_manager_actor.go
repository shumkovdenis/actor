package server

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/maps/treemap"
)

type sessionManagerActor struct {
	sessions    *treemap.Map
	roomManager *actor.PID
}

func newSessionManagerActor() actor.Actor {
	return &sessionManagerActor{
		sessions: treemap.NewWithStringComparator(),
	}
}

func (state *sessionManagerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.roomManager = actor.NewLocalPID("rooms")
	case *CreateSession:
		getRoom := &GetRoom{
			RoomID: msg.RoomID,
		}

		future := state.roomManager.RequestFuture(getRoom, 1*time.Second)
		res, err := future.Result()
		if err != nil {
			err := newErr(ErrCreateSession).Error(err).LogErr()
			ctx.Respond(err)
			return
		}

		if err, ok := res.(*Err); ok {
			err := newErr(ErrCreateSession).Wrap(err).LogErr()
			ctx.Respond(err)
			return
		}

		roomPID := res.(*GetRoomSuccess).RoomPID

		props := actor.FromInstance(newSessionActor(roomPID))
		sessionPID, _ := ctx.SpawnNamed(props, msg.SessionID)

		joinRoom := &JoinRoom{
			SessionPID: sessionPID,
		}

		future = roomPID.RequestFuture(joinRoom, 1*time.Second)
		res, err = future.Result()
		if err != nil {
			err := newErr(ErrCreateSession).Error(err).LogErr()
			ctx.Respond(err)
			return
		}

		if err, ok := res.(*Err); ok {
			err := newErr(ErrCreateSession).Wrap(err).LogErr()
			ctx.Respond(err)
			return
		}

		state.sessions.Put(msg.SessionID, sessionPID)

		success := &CreateSessionSuccess{
			Session: &Session{
				ID:     msg.SessionID,
				RoomID: msg.RoomID,
			},
		}

		ctx.Respond(success)
	case *GetSession:
		pid, ok := state.sessions.Get(msg.SessionID)
		if !ok {
			err := newErr(ErrSessionNotFound).LogErr()
			err = newErr(ErrGetSession).Wrap(err).LogErr()
			ctx.Respond(err)
			return
		}

		success := &GetSessionSuccess{
			SessionPID: pid.(*actor.PID),
		}

		ctx.Respond(success)
	}
}

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
