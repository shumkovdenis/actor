package server

import (
	"go.uber.org/zap"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/shumkovdenis/club/logger"
)

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
			logger.L().Error("create session fail: get room fail",
				zap.Error(err),
			)
			ctx.Respond(&CreateSessionFailed{})
			return
		}

		roomPID, ok := res.(*actor.PID)
		if !ok {
			ctx.Respond(res)
			return
		}

		props := actor.FromInstance(newSessionActor(roomPID))
		sessionPID, _ := ctx.SpawnNamed(props, msg.SessionID)

		joinRoom := &JoinRoom{
			SessionPID: sessionPID,
		}

		future = roomPID.RequestFuture(joinRoom, Timeout)
		res, err = future.Result()
		if err != nil {
			logger.L().Error("create session fail: join room fail",
				zap.Error(err),
			)
			ctx.Respond(&CreateSessionFailed{})
			return
		}

		if _, ok := res.(*JoinedRoom); !ok {
			sessionPID.Stop()
			ctx.Respond(res)
			return
		}

		state.sessions.Put(msg.SessionID, sessionPID)

		ctx.Respond(&Session{msg.SessionID})
	case *GetSession:
		pid, ok := state.sessions.Get(msg.SessionID)
		if !ok {
			ctx.Respond(&SessionNotFound{})
			return
		}

		ctx.Respond(pid)
	}
}
