package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/shumkovdenis/club/messages"
	"github.com/uber-go/zap"
)

type wsActor struct {
	group    *echo.Group
	server   *actor.PID
	upgrader *websocket.Upgrader
}

func newWSActor(group *echo.Group) actor.Actor {
	return &wsActor{
		group:    group,
		upgrader: &websocket.Upgrader{},
	}
}

func (state *wsActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		state.server = ctx.Parent()

		fmt.Println(state.server)

		state.group.POST("", state.connect)
	}
}

func (state *wsActor) connect(c echo.Context) error {
	future := state.server.RequestFuture(&Connect{}, 1*time.Second)
	_, err := future.Result()
	if err != nil {
		log.Error(err.Error())

		return c.String(http.StatusInternalServerError, err.Error())
	}

	ws, err := state.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()

	for {
		cmd := &messages.Command{}
		if err := ws.ReadJSON(cmd); err != nil {
			log.Error(err.Error())

			// ctx.Self().Stop()

			return nil
		}

		log.Debug("Command", zap.String("type", cmd.Type))

		// state.brokerPID.Request(cmd, ctx.Self())

		// if err := ws.WriteJSON(msg); err != nil {
		// 	log.Error(err.Error())
		// }

		// log.Debug("event", zap.String("type", msg.Type))
	}

	return nil
}
