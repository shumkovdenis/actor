package server

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/gorilla/websocket"
	"github.com/uber-go/zap"
)

type wsConnActor struct {
	conn   *websocket.Conn
	reg    Registry
	brkPID *actor.PID
}

func newWSConnActor(conn *websocket.Conn) actor.Actor {
	return &wsConnActor{
		conn: conn,
	}
}

func (state *wsConnActor) SetRegistry(reg Registry) {
	state.reg = reg
}

func (state *wsConnActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromProducer(newBrokerActor).
			WithMiddleware(plugin.Use(RegistryPlugin()))
		pid, err := ctx.SpawnNamed(props, "broker")
		if err != nil {
		}

		state.brkPID = pid

		go state.reader(ctx)

	case *actor.Stopped:
		state.conn.Close()
	case Event:
		evt, err := state.reg.FromMessage(msg)
		if err != nil {
			log.Error(err.Error())

			msg = &Fail{err.Error()}
		}

		log.Debug("Event",
			zap.String("conn", "ws"),
			zap.String("type", evt.Type),
		)

		if err := state.conn.WriteJSON(evt); err != nil {
			log.Error(err.Error())
		}
	}
}

func (state *wsConnActor) reader(ctx actor.Context) {
	defer ctx.Self().Stop()

	var msg interface{}

	for {
		cmd := &command{}
		if err := state.conn.ReadJSON(cmd); err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Error(err.Error())

				return
			}

			log.Error(err.Error())

			msg = &Fail{err.Error()}
		} else {
			log.Debug("Command",
				zap.String("conn", "ws"),
				zap.String("type", cmd.Type),
			)

			msg, err = state.reg.ToMessage(cmd)
			if err != nil {
				log.Error(err.Error())

				msg = &Fail{err.Error()}
			} else {
				future := state.brkPID.RequestFuture(msg, 1*time.Second)
				msg, err = future.Result()
				if err != nil {
					log.Error(err.Error())

					msg = &Fail{err.Error()}
				}
			}
		}

		evt, err := state.reg.FromMessage(msg)
		if err != nil {
			log.Error(err.Error())

			msg = &Fail{err.Error()}
		}

		log.Debug("Event",
			zap.String("conn", "ws"),
			zap.String("type", evt.Type),
		)

		if err := state.conn.WriteJSON(evt); err != nil {
			log.Error(err.Error())
		}
	}
}
