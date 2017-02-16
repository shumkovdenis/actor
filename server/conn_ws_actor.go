package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type wsConnActor struct {
	conn      *websocket.Conn
	conv      Conv
	brokerPID *actor.PID
}

func newWSConnActor(conn *websocket.Conn) actor.Actor {
	return &wsConnActor{
		conn: conn,
		conv: newConv(),
	}
}

func (state *wsConnActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromProducer(newBrokerActor)
		pid, _ := ctx.SpawnNamed(props, "broker")

		state.brokerPID = pid

		go state.reader(ctx)
	case Event:
		log.Debug("event",
			zap.String("conn", "ws"),
			zap.String("type", msg.Event()),
		)

		evt, err := state.conv.FromMessage(msg)
		if err != nil {
			log.Error(err.Error())

			return
		}

		if err := state.conn.WriteJSON(evt); err != nil {
			log.Error(err.Error())
		}
	}
}

func (state *wsConnActor) reader(ctx actor.Context) {
	defer ctx.Self().Stop()

	for {
		cmd := &command{}
		if err := state.conn.ReadJSON(cmd); err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Error(err.Error())
				return
			}

			// err := newErr(ErrReadJSON).Error(err).LogErr()
			// ctx.Self().Tell(err)
			continue
		}

		log.Debug("command",
			zap.String("conn", "ws"),
			zap.String("type", cmd.Type),
		)

		msg, err := state.conv.ToMessage(cmd)
		if err != nil {
			// err := newErr(ErrToMessage).Error(err).LogErr()
			// ctx.Self().Tell(err)
			continue
		}

		state.brokerPID.Request(msg, ctx.Self())
	}
}
