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
		}

		state.sesMngPID = pid

		state.grp.POST("/sessions", state.createSession)
	}
}

func (state *apiActor) createSession(c echo.Context) error {
	session := &struct {
		RoomID string `json:"room_id"`
	}{}

	if err := c.Bind(session); err != nil {
		return err
	}

	future := state.sesMngPID.RequestFuture(&CreateSession{}, 1*time.Second)
	res, err := future.Result()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
