package server

import (
	"net/http"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/labstack/echo"
)

type httpConnActor struct {
	grp    *echo.Group
	brk    Broker
	msgs   *arraylist.List
	brkPID *actor.PID
}

func newHTTPConnActor(group *echo.Group) actor.Actor {
	return &httpConnActor{
		grp:  group,
		brk:  newBroker(),
		msgs: arraylist.New(),
	}
}

func (state *httpConnActor) Receive(ctx actor.Context) {
	// switch msg := ctx.Message().(type) {
	// case *actor.Started:
	// 	props := actor.FromInstance(newBrokerActor())
	// 	pid, err := ctx.SpawnNamed(props, "broker")
	// 	if err != nil {
	// 		log.Error(err.Error())
	// 	}

	// 	state.brkPID = pid

	// 	state.grp.POST("/push", state.push)
	// 	state.grp.POST("/pull", state.pull)
	// case Event:
	// 	evt, err := state.reg.FromMessage(msg)
	// 	if err != nil {
	// 		log.Error(err.Error())

	// 		return
	// 	}

	// 	sub := state.brk.Contains(evt.Type)

	// 	log.Debug("Event",
	// 		zap.String("conn", "ws"),
	// 		zap.String("type", evt.Type),
	// 		zap.Bool("sub", sub),
	// 	)

	// 	if sub {
	// 		state.msgs.Add(evt)
	// 	}
	// }
}

func (state *httpConnActor) push(c echo.Context) error {
	// cmd := &command{}
	// if err := c.Bind(cmd); err != nil {
	// 	return err
	// }

	// log.Debug("Command",
	// 	zap.String("conn", "http"),
	// 	zap.String("type", cmd.Type),
	// )

	// msg, err := state.reg.ToMessage(cmd)
	// if err != nil {
	// 	return err
	// }

	// future := state.brkPID.RequestFuture(msg, 1*time.Second)
	// msg, err = future.Result()
	// if err != nil {
	// 	return err
	// }

	// evt, err := state.reg.FromMessage(msg)
	// if err != nil {
	// 	return err
	// }

	// sub := state.brk.Contains(evt.Type)

	// log.Debug("Event",
	// 	zap.String("conn", "http"),
	// 	zap.String("type", evt.Type),
	// 	zap.Bool("sub", sub),
	// )

	// if !sub {
	// 	return errors.New("no subscription event: " + evt.Type)
	// }

	// return c.JSON(http.StatusOK, evt)
	return nil
}

func (state *httpConnActor) pull(c echo.Context) error {
	msgs := state.msgs.Values()

	state.msgs.Clear()

	return c.JSON(http.StatusOK, msgs)
}
