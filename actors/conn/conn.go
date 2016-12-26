package conn

import (
	"log"

	"github.com/AsynkronIT/gam/actor"
	"github.com/gorilla/websocket"
	"github.com/shumkovdenis/actor/actors/broker"
	"github.com/shumkovdenis/actor/messages"
)

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
		state.reader(ctx)
	case *messages.Event:
		err := state.ws.WriteJSON(msg)
		if err != nil {
			log.Fatalf("write error: %s", err)
		}
	}
}

func (state *connActor) reader(ctx actor.Context) {
	defer state.ws.Close()
	for {
		cmd := &messages.Command{}
		if err := state.ws.ReadJSON(cmd); err != nil {
			log.Fatalf("read error: %s\n", err)
		}
		state.brokerPID.Request(cmd, ctx.Self())
	}
}
