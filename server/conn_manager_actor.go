package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type connManagerActor struct {
	grp *echo.Group
	upg *websocket.Upgrader
}

func newConnManagerActor(group *echo.Group) actor.Actor {
	return &connManagerActor{
		grp: group,
		upg: &websocket.Upgrader{},
	}
}

func (state *connManagerActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		// state.grp.GET("/http", state.http(ctx))
		state.grp.GET("/ws", state.ws(ctx))

		log.Info("Conn manager started")
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
		conn, err := state.upg.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Error(err.Error())

			return err
		}

		id := uuid.NewV4().String()

		props := actor.FromInstance(newWSConnActor(conn))
		_, err = ctx.SpawnNamed(props, id)
		if err != nil {
			log.Error(err.Error())

			return err
		}

		return nil
	}
}
