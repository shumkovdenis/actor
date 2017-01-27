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
		_, err := ctx.SpawnNamed(props, "api")
		if err != nil {
		}

		props = actor.FromInstance(newConnManagerActor(e.Group("/conn")))
		_, err = ctx.SpawnNamed(props, "conns")
		if err != nil {
		}

		go func() {
			e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Port)))
		}()
	}
}
