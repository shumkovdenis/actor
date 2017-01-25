package server

import "github.com/AsynkronIT/protoactor-go/actor"

type Server interface {
	Registry() Registry
}

type serverActor struct {
	reg Registry
}

func NewServerActor() actor.Actor {
	return &serverActor{
		reg: newRegistry(),
	}
}

func (state *serverActor) Registry() Registry {
	return state.reg
}

func (state *serverActor) Receive(ctx actor.Context) {

}
