package serv

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	uuid "github.com/satori/go.uuid"
	"github.com/shumkovdenis/club/actors/conn"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
	"github.com/shumkovdenis/club/plugins"
)

var log = logger.Get()

type serverActor struct {
	brokerList *plugins.List
	ratesSubs  *plugins.Subs
}

var (
	upgrader = websocket.Upgrader{}
)

func New(brokerList *plugins.List) actor.Actor {
	return &serverActor{
		brokerList: brokerList,
		ratesSubs:  plugins.NewSubs("event.rates.change"),
	}
}

func (state *serverActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		state.start(ctx)

		log.Info("Server actor started")
	}
}

func (state *serverActor) start(ctx actor.Context) {
	conf := config.Server()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Static("/", conf.PublicPath)
	e.GET("/ws", state.handler(ctx))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Port)))
}

func (state *serverActor) handler(ctx actor.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		id := uuid.NewV4().String()
		props := actor.FromInstance(conn.New(state.brokerList, state.ratesSubs, ws))
		ctx.SpawnNamed(props, id)

		return nil
	}
}
