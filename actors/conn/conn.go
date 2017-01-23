package conn

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gorilla/websocket"
	"github.com/shumkovdenis/club/actors/broker"
	"github.com/shumkovdenis/club/logger"
	"github.com/shumkovdenis/club/messages"
	"github.com/uber-go/zap"
)

var log = logger.Get()

type connActor struct {
	ws        *websocket.Conn
	brokerPID *actor.PID
}

func New(ws *websocket.Conn) actor.Actor {
	return &connActor{
		ws: ws,
	}
}

func (state *connActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromProducer(broker.NewActor)
		state.brokerPID = ctx.Spawn(props)
		go state.reader(ctx)
	case *messages.Event:
		err := state.ws.WriteJSON(msg)
		if err != nil {
			log.Error("write error", zap.Error(err))
		}

		log.Debug("event", zap.String("type", msg.Type))
	}
}

func (state *connActor) reader(ctx actor.Context) {
	defer state.ws.Close()
	for {
		cmd := &messages.Command{}
		if err := state.ws.ReadJSON(cmd); err != nil {
			log.Error("read error", zap.Error(err))
			ctx.Self().Stop()
			return
		}

		log.Debug("command", zap.String("type", cmd.Type))

		state.brokerPID.Request(cmd, ctx.Self())
	}
}
