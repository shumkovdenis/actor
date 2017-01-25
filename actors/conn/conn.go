package conn

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/gorilla/websocket"
	"github.com/shumkovdenis/club/actors/broker"
	"github.com/shumkovdenis/club/logger"
	"github.com/shumkovdenis/club/messages"
	"github.com/shumkovdenis/club/plugins"
	"github.com/uber-go/zap"
)

var log = logger.Get()

type connActor struct {
	brokerList *plugins.List
	ratesSubs  *plugins.Subs
	ws         *websocket.Conn
	brokerPID  *actor.PID
}

func New(brokerList *plugins.List, ratesSubs *plugins.Subs,
	ws *websocket.Conn) actor.Actor {
	return &connActor{
		brokerList: brokerList,
		ws:         ws,
	}
}

func (state *connActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromProducer(broker.New).
			WithMiddleware(
				plugin.Use(state.brokerList),
				plugin.Use(state.ratesSubs),
			)
		state.brokerPID = ctx.Spawn(props)

		go state.reader(ctx)

		log.Info("Conn actor started")
	case *messages.Event:
		if err := state.ws.WriteJSON(msg); err != nil {
			log.Error(err.Error())
		}

		log.Debug("event", zap.String("type", msg.Type))
	}
}

func (state *connActor) reader(ctx actor.Context) {
	defer state.ws.Close()
	for {
		cmd := &messages.Command{}
		if err := state.ws.ReadJSON(cmd); err != nil {
			log.Error(err.Error())

			ctx.Self().Stop()

			return
		}

		log.Debug("Command", zap.String("type", cmd.Type))

		state.brokerPID.Request(cmd, ctx.Self())
	}
}
