package server

import (
	"net/http"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/labstack/echo"
)

type apiActor struct {
	grp       *echo.Group
	sesMngPID *actor.PID
	roomMng   *actor.PID
}

func newAPIActor(group *echo.Group) actor.Actor {
	return &apiActor{
		grp: group,
	}
}

func (state *apiActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromProducer(newSessionManagerActor)
		pid, err := ctx.SpawnNamed(props, "sessions")
		if err != nil {
			log.Error(err.Error())
		} else {
			state.sesMngPID = pid

			state.grp.POST("/sessions", state.createSession)
		}

		props = actor.FromProducer(newRoomManagerActor)
		pid, err = ctx.SpawnNamed(props, "rooms")
		if err != nil {
			log.Error(err.Error())
		} else {
			state.roomMng = pid

			state.grp.POST("/rooms", state.createRoom)
		}
	}
}

func (state *apiActor) createSession(c echo.Context) error {
	conf := newSessionConf()

	if err := c.Bind(conf); err != nil {
		return c.JSON(http.StatusBadRequest, ErrParsing)
	}

	future := state.sesMngPID.RequestFuture(&CreateSession{conf}, 1*time.Second)
	res, err := future.Result()
	if err != nil {
		return err
	}

	if fail, ok := res.(CreateSessionFail); ok {
		return c.JSON(http.StatusBadRequest, &ClientError{fail.Error()})
	}

	session := res.(CreateSessionSuccess).Session

	future = state.roomMng.RequestFuture(&JoinRoom{session}, 1*time.Second)
	res, err = future.Result()
	if err != nil {
		return err
	}

	if fail, ok := res.(JoinRoomFail); ok {
		return c.JSON(http.StatusBadRequest, &ClientError{fail.Error()})
	}

	return c.JSON(http.StatusOK, session)
}

func (state *apiActor) createRoom(c echo.Context) error {
	future := state.roomMng.RequestFuture(&CreateRoom{}, 1*time.Second)
	res, err := future.Result()
	if err != nil {
		return err
	}

	if fail, ok := res.(CreateRoomFail); ok {
		return c.JSON(http.StatusBadRequest, &ClientError{fail.Error()})
	}

	room := res.(*CreateRoomSuccess).Room

	return c.JSON(http.StatusOK, room)
}
