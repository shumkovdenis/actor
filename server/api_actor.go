package server

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"github.com/shumkovdenis/club/server/core"
)

const (
	Timeout = 1 * time.Second
)

type apiActor struct {
	group             *echo.Group
	roomManagerPID    *actor.PID
	sessionManagerPID *actor.PID
}

func newAPIActor(group *echo.Group) actor.Actor {
	return &apiActor{
		group: group,
	}
}

func (state *apiActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		state.roomManagerPID = actor.NewLocalPID("rooms")
		state.sessionManagerPID = actor.NewLocalPID("sessions")

		state.group.POST("/rooms", state.createRoom)
		state.group.POST("/sessions", state.createSession)
	}
}

func (state *apiActor) createRoom(c echo.Context) error {
	createRoom := &CreateRoom{
		RoomID: uuid.NewV4().String(),
	}

	future := state.roomManagerPID.RequestFuture(createRoom, Timeout)
	res, err := future.Result()
	if err != nil {
		log.Error("api create room failed: create room failed",
			zap.Error(err),
		)
		return c.NoContent(http.StatusInternalServerError)
	}

	room, ok := res.(*Room)
	if !ok {
		if fail, ok := res.(core.Fail); ok {
			return c.JSON(http.StatusBadRequest, failResp(fail))
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := &struct {
		ID string `json:"id"`
	}{
		ID: room.ID,
	}

	return c.JSON(http.StatusOK, resp)
}

func (state *apiActor) createSession(c echo.Context) error {
	req := &struct {
		RoomID string `json:"room_id"`
	}{}

	if err := c.Bind(req); err != nil {
		log.Error("api create session failed: parse json failed",
			zap.Error(err),
		)
		return c.JSON(http.StatusBadRequest, failResp(&ParseJSONFailed{}))
	}

	createSession := &CreateSession{
		SessionID: uuid.NewV4().String(),
		RoomID:    req.RoomID,
	}

	future := state.sessionManagerPID.RequestFuture(createSession, Timeout)
	res, err := future.Result()
	if err != nil {
		log.Error("api create session failed: create session failed",
			zap.Error(err),
		)
		return c.NoContent(http.StatusInternalServerError)
	}

	session, ok := res.(*Session)
	if !ok {
		if fail, ok := res.(core.Fail); ok {
			return c.JSON(http.StatusBadRequest, failResp(fail))
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := &struct {
		ID string `json:"id"`
	}{
		ID: session.ID,
	}

	return c.JSON(http.StatusOK, resp)
}
