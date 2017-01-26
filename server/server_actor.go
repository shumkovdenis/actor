package server

import (
	"fmt"
	"net/http"
	"time"

	"errors"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"github.com/shumkovdenis/club/config"
	"github.com/uber-go/zap"
)

type Server interface {
	Registry() Registry
}

type serverActor struct {
	registry          Registry
	context           actor.Context
	cons              *treemap.Map
	upgrader          *websocket.Upgrader
	sessionManagerPID *actor.PID
}

func newServerActor() actor.Actor {
	return &serverActor{
		registry: newRegistry(),
		cons:     treemap.NewWithStringComparator(),
		upgrader: &websocket.Upgrader{},
	}
}

func (state *serverActor) Registry() Registry {
	return state.registry
}

func (state *serverActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		state.context = ctx

		props := actor.FromProducer(newSessionManagerActor)
		pid, err := ctx.SpawnNamed(props, "sessions")
		if err != nil {
		}
		state.sessionManagerPID = pid

		conf := config.Server()

		e := echo.New()
		e.Static("/", conf.PublicPath)
		e.GET("/ws", state.ws)

		http := e.Group("/http")
		http.POST("/connect", state.httpConnect)
		http.POST("/disconnect/:id", state.httpDisconnect)
		http.POST("/command/:id", state.httpCommand)

		api := e.Group("/api")
		api.POST("/sessions", state.apiSessions)

		go func() {
			e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Port)))
		}()
	}
}

func (state *serverActor) connect() (string, error) {
	id := uuid.NewV4().String()

	props := actor.FromProducer(newConnActor).
		WithMiddleware(plugin.Use(RegistryPlugin()))
	pid, err := state.context.SpawnNamed(props, id)
	if err != nil {
		return "", err
	}

	state.cons.Put(id, pid)

	return id, nil
}

func (state *serverActor) disconnect(id string) error {
	pid, ok := state.cons.Get(id)
	if !ok {
		return errors.New("Connection not found")
	}
	pid.(*actor.PID).Stop()

	state.cons.Remove(id)

	return nil
}

func (state *serverActor) proccess(id string, cmd *command) (*event, error) {
	pid, ok := state.cons.Get(id)
	if !ok {
		return nil, errors.New("Connection not found")
	}

	log.Debug("Command", zap.String("type", cmd.Type))

	msg, err := state.registry.ToMessage(cmd)
	if err != nil {
		return nil, err
	}

	future := pid.(*actor.PID).RequestFuture(msg, 1*time.Second)
	res, err := future.Result()
	if err != nil {
		return nil, errors.New("Failure command processing")
	}

	evt, err := state.registry.FromMessage(res)
	if err != nil {
		return nil, err
	}

	log.Debug("Event", zap.String("type", evt.Type))

	return evt, nil
}

func (state *serverActor) ws(c echo.Context) error {
	id, err := state.connect()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	ws, err := state.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	defer ws.Close()

	for {
		cmd := &command{}
		if err := ws.ReadJSON(cmd); err != nil {
			state.disconnect(id)

			return nil
		}

		evt, err := state.proccess(id, cmd)
		if err != nil {
		}

		if err := ws.WriteJSON(evt); err != nil {
		}
	}

	return nil
}

func (state *serverActor) httpConnect(c echo.Context) error {
	id, err := state.connect()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	resp := &struct {
		ConnectionID string `json:"connection_id"`
	}{
		ConnectionID: id,
	}

	return c.JSON(http.StatusOK, resp)
}

func (state *serverActor) httpDisconnect(c echo.Context) error {
	id := c.Param("id")

	if err := state.disconnect(id); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	resp := &struct{}{}

	return c.JSON(http.StatusOK, resp)
}

func (state *serverActor) httpCommand(c echo.Context) error {
	id := c.Param("id")

	cmd := &command{}
	if err := c.Bind(cmd); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	evt, err := state.proccess(id, cmd)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, evt)
}

func (state *serverActor) apiSessions(c echo.Context) error {
	session := &struct {
		RoomID string `json:"room_id"`
	}{}

	if err := c.Bind(session); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	future := state.sessionManagerPID.RequestFuture(&CreateSession{}, 1*time.Second)
	res, err := future.Result()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
