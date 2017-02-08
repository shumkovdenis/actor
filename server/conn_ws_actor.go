package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/gorilla/websocket"
	"github.com/uber-go/zap"
)

type wsConnActor struct {
	conn   *websocket.Conn
	brk    Broker
	reg    Registry
	brkPID *actor.PID
}

func newWSConnActor(conn *websocket.Conn) actor.Actor {
	return &wsConnActor{
		conn: conn,
		brk:  newBroker(),
	}
}

func (state *wsConnActor) Init(reg Registry) {
	state.reg = reg
}

func (state *wsConnActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromInstance(newBrokerActor(state.brk)).
			WithMiddleware(plugin.Use(RegistryPlugin()))
		pid, err := ctx.SpawnNamed(props, "broker")
		if err != nil {
			log.Error(err.Error())
		}

		state.brkPID = pid

		go state.reader(ctx)

	case *actor.Stopped:
		state.conn.Close()
	case Event:
		evt, err := state.reg.FromMessage(msg)
		if err != nil {
			log.Error(err.Error())

			return
		}

		sub := state.brk.Contains(evt.Type)

		log.Debug("Event",
			zap.String("conn", "ws"),
			zap.String("type", evt.Type),
			zap.Bool("sub", sub),
		)

		if sub {
			if err := state.conn.WriteJSON(evt); err != nil {
				log.Error(err.Error())
			}
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

			log.Error(err.Error())

			// ctx.Self().Tell(&Fail{err.Error()})

			continue
		}

		log.Debug("Command",
			zap.String("conn", "ws"),
			zap.String("type", cmd.Type),
		)

		msg, err := state.reg.ToMessage(cmd)
		if err != nil {
			log.Error(err.Error())

			// ctx.Self().Tell(&Fail{err.Error()})

			continue
		}

		state.brkPID.Request(msg, ctx.Self())
	}
}
