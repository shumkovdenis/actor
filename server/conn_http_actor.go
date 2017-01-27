package server

import (
	"net/http"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/labstack/echo"
)

type httpConnActor struct {
	grp    *echo.Group
	reg    Registry
	brkPID *actor.PID
}

func newHTTPConnActor(group *echo.Group) actor.Actor {
	return &httpConnActor{
		grp: group,
	}
}

func (state *httpConnActor) SetRegistry(reg Registry) {
	state.reg = reg
}

func (state *httpConnActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromProducer(newBrokerActor).
			WithMiddleware(plugin.Use(RegistryPlugin()))
		pid, err := ctx.SpawnNamed(props, "broker")
		if err != nil {
		}

		state.brkPID = pid

		state.grp.POST("/push", state.push)
		state.grp.POST("/pull", state.pull)
	case Event:
	}
}

func (state *httpConnActor) push(c echo.Context) error {
	cmd := &command{}
	if err := c.Bind(cmd); err != nil {
		return err
	}

	msg, err := state.reg.ToMessage(cmd)
	if err != nil {
		return err
	}

	future := state.brkPID.RequestFuture(msg, 1*time.Second)
	msg, err = future.Result()
	if err != nil {
		return err
	}

	evt, err := state.reg.FromMessage(msg)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, evt)
}

func (*httpConnActor) pull(c echo.Context) error {
	return nil
}
