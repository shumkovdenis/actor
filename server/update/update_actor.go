package update

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

type updateActor struct{}

func NewActor() actor.Actor {
	return &updateActor{}
}

func (state *updateActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromProducer(newAutoUpdateActor)
		ctx.SpawnNamed(props, "auto")

		ctx.SetBehavior(state.started)
	}
}

func (state *updateActor) started(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *Check:
		ctx.SetBehavior(state.checking)
		res := check()
		ctx.Respond(res)
		ctx.SetBehavior(state.started)
	case *Download:
		ctx.SetBehavior(state.downloading)
		ch := download()
		for res := range ch {
			ctx.Respond(res)
		}
		ctx.SetBehavior(state.started)
	case *Install:
		ctx.SetBehavior(state.installing)
		res := install()
		ctx.Respond(res)
		ctx.SetBehavior(state.started)
	}
}

func (state *updateActor) checking(ctx actor.Context) {
	switch ctx.Message().(type) {
	case Message:
		ctx.Respond(&Checking{})
	}
}

func (state *updateActor) downloading(ctx actor.Context) {
	switch ctx.Message().(type) {
	case Message:
		ctx.Respond(&Downloading{})
	}
}

func (state *updateActor) installing(ctx actor.Context) {
	switch ctx.Message().(type) {
	case Message:
		ctx.Respond(&Installing{})
	}
}
