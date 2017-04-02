package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gorilla/websocket"
	"github.com/shumkovdenis/club/logger"
	"github.com/shumkovdenis/club/server/core"
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
	case *actor.Stopped:
		if err := state.conn.Close(); err != nil {
			logger.L().Error("close websocket connection failed",
				zap.Error(err),
			)
		}
	case core.Event:
		l := logger.L().With(
			zap.String("conn", "ws"),
			zap.String("type", msg.Event()),
		)

		if code, ok := msg.(core.Code); ok {
			l = l.With(zap.String("code", code.Code()))
		}

		l.Debug("event")

		evt, err := state.conv.FromMessage(msg)
		if err != nil {
			logger.L().Error("conv from message failed",
				zap.Error(err),
			)
			return
		}

		if err := state.conn.WriteJSON(evt); err != nil {
			logger.L().Error("write websocket failed",
				zap.Error(err),
			)
		}
	}
}

func (state *wsConnActor) reader(ctx actor.Context) {
	defer ctx.Self().Stop()

	for {
		cmd := &command{}
		if err := state.conn.ReadJSON(cmd); err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				logger.L().Warn("close websocket",
					zap.Error(err),
				)
				return
			}
			logger.L().Error("read websocket failed",
				zap.Error(err),
			)
			ctx.Self().Tell(&ConnReadFailed{})
			continue
		}

		logger.L().Debug("command",
			zap.String("conn", "ws"),
			zap.String("type", cmd.Type),
		)

		msg, err := state.conv.ToMessage(cmd)
		if err != nil {
			logger.L().Error("conv to message failed",
				zap.Error(err),
			)
			ctx.Self().Tell(&ConnReadFailed{})
			continue
		}

		state.brokerPID.Request(msg, ctx.Self())
	}
}
