package server

import (
	"net/http"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type apiActor struct {
	grp            *echo.Group
	roomManager    *actor.PID
	sessionManager *actor.PID
}

func newAPIActor(group *echo.Group) actor.Actor {
	return &apiActor{
		grp: group,
	}
}

func (state *apiActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		state.roomManager = actor.NewLocalPID("rooms")
		state.sessionManager = actor.NewLocalPID("sessions")

		state.grp.POST("/rooms", state.createRoom)
		state.grp.POST("/sessions", state.createSession)
	}
}

func (state *apiActor) createRoom(c echo.Context) error {
	createRoom := &CreateRoom{
		RoomID: uuid.NewV4().String(),
	}

	future := state.roomManager.RequestFuture(createRoom, 1*time.Second)
	res, err := future.Result()
	if err != nil {
		err := newErrorWrap(ErrCreateRoom, err)
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err, ok := res.(*Error); ok {
		err := newErrorWrap(ErrCreateRoom, err)
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, err)
	}

	success := res.(*CreateRoomSuccess)

	return c.JSON(http.StatusOK, success.Room)
}

func (state *apiActor) createSession(c echo.Context) error {
	createSession := &CreateSession{
		SessionID: uuid.NewV4().String(),
	}

	if err := c.Bind(createSession); err != nil {
		err := newErrorWrap(ErrCreateSession, err)
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, err)
	}

	future := state.sessionManager.RequestFuture(createSession, 1*time.Second)
	res, err := future.Result()
	if err != nil {
		err := newErrorWrap(ErrCreateSession, err)
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err, ok := res.(*Error); ok {
		err := newErrorWrap(ErrCreateSession, err)
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, err)
	}

	success := res.(*CreateSessionSuccess)

	return c.JSON(http.StatusOK, success.Session)
}
