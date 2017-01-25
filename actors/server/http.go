package server

import (
	"net/http"
	"time"

	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/labstack/echo"
	"github.com/shumkovdenis/club/messages"
)

type httpActor struct {
	group  *echo.Group
	server *actor.PID
}

func newHTTPActor(group *echo.Group) actor.Actor {
	return &httpActor{
		group: group,
	}
}

func (state *httpActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		state.server = ctx.Parent()

		fmt.Println(state.server)

		state.group.POST("/connect", state.connect)
		state.group.POST("/disconnect/:id", state.disconnect)
		state.group.POST("/message/:id", state.message)
	}
}

func (state *httpActor) connect(c echo.Context) error {
	future := state.server.RequestFuture(&Connect{}, 1*time.Second)
	msg, err := future.Result()
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, msg)
}

func (state *httpActor) disconnect(c echo.Context) error {
	id := c.Param("id")

	future := state.server.RequestFuture(&Disconnect{id}, 1*time.Second)
	msg, err := future.Result()
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, msg)
}

func (state *httpActor) message(c echo.Context) error {
	id := c.Param("id")

	cmd := &messages.Command{}
	if err := c.Bind(cmd); err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, err.Error())
	}

	future := state.server.RequestFuture(&Message{id, cmd}, 1*time.Second)
	evt, err := future.Result()
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, evt)
}
