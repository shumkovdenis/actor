package server

import (
	"net/http"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"github.com/shumkovdenis/club/logger"
	"go.uber.org/zap"
)

type connManagerActor struct {
	group    *echo.Group
	upgrader *websocket.Upgrader
}

func newConnManagerActor(group *echo.Group) actor.Actor {
	return &connManagerActor{
		group: group,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (state *connManagerActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		// state.grp.GET("/http", state.http(ctx))
		state.group.GET("/ws", state.ws(ctx))
	}
}

// func (state *connManagerActor) http(ctx actor.Context) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		id := uuid.NewV4().String()

// 		group := state.grp.Group("/" + id)

// 		props := actor.FromInstance(newHTTPConnActor(group)).
// 			WithMiddleware(plugin.Use(RegistryPlugin()))
// 		_, err := ctx.SpawnNamed(props, id)
// 		if err != nil {
// 			return err
// 		}

// 		resp := &struct {
// 			ConnectionID string `json:"connection_id"`
// 		}{
// 			ConnectionID: id,
// 		}

// 		return c.JSON(http.StatusOK, resp)
// 	}
// }

func (state *connManagerActor) ws(ctx actor.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := state.upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			logger.L().Error("create websocket connection failed",
				zap.Error(err),
			)
			return err
		}

		id := uuid.NewV4().String()

		props := actor.FromInstance(newWSConnActor(conn))
		ctx.SpawnNamed(props, id)

		return nil
	}
}
