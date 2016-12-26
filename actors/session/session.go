package session

import (
	"strings"

	"github.com/AsynkronIT/gam/actor"
	uuid "github.com/satori/go.uuid"
	"github.com/shumkovdenis/actor/actors/client"
)

type Login struct {
	Client string `json:"client"`
}

type LoginSuccess struct {
	Client string `json:"client"`
}

type LoginFail struct {
}

type sessionActor struct {
	path      string
	clientPID *actor.PID
}

func NewActor() actor.Actor {
	return &sessionActor{}
}

func (state *sessionActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Login:
		var id = msg.Client
		if len(strings.TrimSpace(id)) == 0 {
			state.path = "app"
			props := actor.FromProducer(client.NewActor)
			id = uuid.NewV4().String()
			state.clientPID = actor.SpawnNamed(props, "/clients/"+id)
		} else {
			state.path = "web"
			state.clientPID = actor.NewLocalPID("/clients/" + id)
		}
		state.clientPID.Request(&client.Join{Client: id}, ctx.Self())
	case *client.Joined:
		ctx.Parent().Tell(&LoginSuccess{msg.Client})
		switch state.path {
		case "app":
			ctx.Become(state.App)
		case "web":
			ctx.Become(state.Web)
		}
	default:
		state.clientPID.Request(msg, ctx.Parent())
	}
}

func (state *sessionActor) App(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	}
}

func (state *sessionActor) Web(ctx actor.Context) {

}
