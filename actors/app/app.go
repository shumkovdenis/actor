package app

import "github.com/AsynkronIT/gam/actor"

type Fail struct {
	Message string `json:"message"`
}

type appActor struct {
}

func NewActor() actor.Actor {
	return &appActor{}
}

func (state *appActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
	}
}
