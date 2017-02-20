package update

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/config"
)

type updateActor struct {
}

func NewActor() actor.Actor {
	return &updateActor{}
}

func (state *updateActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		conf := config.UpdateServer()

		if conf.AutoUpdate {
		}

		ctx.SetBehavior(state.started)
	}
}

func (state *updateActor) started(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *Check:
		ctx.SetBehavior(state.checking)

		// actors.Process(check, ctx.Respond)

		ctx.SetBehavior(state.started)
	case *Download:
		ctx.SetBehavior(state.downloading)

		// actors.Process(download, ctx.Respond)

		ctx.SetBehavior(state.started)
	case *Install:
		ctx.SetBehavior(state.installing)

		// actors.Process(install, ctx.Respond)

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
