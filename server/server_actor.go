package server

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/labstack/echo"
	"github.com/shumkovdenis/club/config"
)

type Server interface {
	Registry() Registry
}

type serverActor struct {
	registry Registry
}

func newServerActor() actor.Actor {
	return &serverActor{
		registry: newRegistry(),
	}
}

func (state *serverActor) Registry() Registry {
	return state.registry
}

func (state *serverActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		conf := config.Server()

		e := echo.New()

		props := actor.FromInstance(newAPIActor(e.Group("/api")))
		ctx.SpawnNamed(props, "api")

		props = actor.FromInstance(newConnManagerActor(e.Group("/conn")))
		ctx.SpawnNamed(props, "conns")

		// go func() {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Port)))
		// }()
	}
}
