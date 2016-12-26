package server

import (
	"github.com/AsynkronIT/gam/actor"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/shumkovdenis/actor/actors/conn"
)

type serverActor struct {
}

var (
	upgrader = websocket.Upgrader{}
)

func New() actor.Actor {
	return &serverActor{}
}

func (state *serverActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		state.start(ctx)
	}
}

func (state *serverActor) start(ctx actor.Context) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/ws", state.handler(ctx))
	e.Logger.Fatal(e.Start(":8282"))
}

func (state *serverActor) handler(ctx actor.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		props := actor.FromInstance(conn.New(ws))
		ctx.Spawn(props)

		return nil
	}
}
