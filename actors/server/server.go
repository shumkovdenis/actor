package server

import (
	"fmt"

	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
	"github.com/shumkovdenis/club/messages"
)

var log = logger.Get()

func init() {

}

type Command struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Connect struct {
}

type Connected struct {
	ConnectionID string `json:"connection_id"`
}

type Disconnect struct {
	ConnectionID string
}

type Disconnected struct {
}

type Message struct {
	ConnectionID string
	Command      *messages.Command
}

type Fail struct {
	Message string `json:"message"`
}

type Server interface {
	Registry() Registry
}

type serverActor struct {
	cons *treemap.Map
}

func NewServerActor() actor.Actor {
	return &serverActor{
		cons: treemap.NewWithStringComparator(),
	}
}

func (state *serverActor) Server() {}

func (state *serverActor) Receive(ctx actor.Context) {
	log.Debug(fmt.Sprintf("%v", ctx.Message()))
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		conf := config.Server()

		e := echo.New()
		e.Static("/", conf.PublicPath)

		props := actor.FromInstance(newHTTPActor(e.Group("/http")))
		_, err := ctx.SpawnNamed(props, "htpp")
		if err != nil {
			log.Error(err.Error())
		}

		props = actor.FromInstance(newWSActor(e.Group("/ws")))
		_, err = ctx.SpawnNamed(props, "ws")
		if err != nil {
			log.Error(err.Error())
		}

		go func() {
			e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Port)))
		}()

		log.Info("Server actor started")
	case *Connect:
		id := uuid.NewV4().String()

		props := actor.FromProducer(newBrokerActor)
		pid, err := ctx.SpawnNamed(props, id)
		if err != nil {
			log.Error(err.Error())
		}

		state.cons.Put(id, pid)

		ctx.Respond(&Connected{id})
	case *Disconnect:
		pid, ok := state.cons.Get(msg.ConnectionID)
		if !ok {
			ctx.Respond(&Fail{"Connection not found"})

			return
		}

		state.cons.Remove(msg.ConnectionID)

		pid.(*actor.PID).Stop()

		ctx.Respond(&Disconnected{})
	case *Message:
		pid, ok := state.cons.Get(msg.ConnectionID)
		if !ok {
			ctx.Respond(&Fail{"Connection not found"})

			return
		}

		future := pid.(*actor.PID).RequestFuture(msg.Command, 1*time.Second)
		res, err := future.Result()
		if err != nil {
			log.Error(err.Error())

			ctx.Respond(&Fail{"Failure command processing"})

			return
		}

		ctx.Respond(res)
	}
}
