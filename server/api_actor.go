package server

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

const (
	ParseJSONFail = "parse_json_fail"

	Timeout = 1 * time.Second
)

type apiActor struct {
	grp               *echo.Group
	roomManagerPID    *actor.PID
	sessionManagerPID *actor.PID
}

func newAPIActor(group *echo.Group) actor.Actor {
	return &apiActor{
		grp: group,
	}
}

func (state *apiActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		state.roomManagerPID = actor.NewLocalPID("rooms")
		state.sessionManagerPID = actor.NewLocalPID("sessions")

		state.grp.POST("/rooms", state.createRoom)
		state.grp.POST("/sessions", state.createSession)
	}
}

func (state *apiActor) createRoom(c echo.Context) error {
	createRoom := &CreateRoom{
		RoomID: uuid.NewV4().String(),
	}

	future := state.roomManagerPID.RequestFuture(createRoom, Timeout)
	res, err := future.Result()
	if err != nil {
		log.Error("api create room fail: create room fail",
			zap.Error(err),
		)
		return c.NoContent(http.StatusInternalServerError)
	}

	if fail, ok := res.(Fail); ok {
		return c.JSON(http.StatusBadRequest, fail)
	}

	return c.JSON(http.StatusOK, res.(*CreateRoomSuccess).Room)
}

func (state *apiActor) createSession(c echo.Context) error {
	createSession := &CreateSession{
		SessionID: uuid.NewV4().String(),
	}

	if err := c.Bind(createSession); err != nil {
		log.Error("api create session fail: parse json fail",
			zap.Error(err),
		)
		return c.JSON(http.StatusBadRequest, newFail(ParseJSONFail))
	}

	future := state.sessionManagerPID.RequestFuture(createSession, Timeout)
	res, err := future.Result()
	if err != nil {
		log.Error("api create session fail: create session fail",
			zap.Error(err),
		)
		return c.NoContent(http.StatusInternalServerError)
	}

	if fail, ok := res.(Fail); ok {
		return c.JSON(http.StatusBadRequest, fail)
	}

	return c.JSON(http.StatusOK, res.(*CreateSessionSuccess).Session)
}
