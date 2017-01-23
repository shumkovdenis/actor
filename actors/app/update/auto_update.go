package update

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/actors/group"
	"github.com/shumkovdenis/club/config"
)

type autoUpdateActor struct {
	updater  *actor.PID
	listener *actor.PID
}

func newAutoUpdate(updater *actor.PID, listener *actor.PID) actor.Actor {
	return &autoUpdateActor{
		updater:  updater,
		listener: listener,
	}
}

func (state *autoUpdateActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		state.listener.Request(&group.Use{
			Producer: ctx.Self(),
			Types: []interface{}{
				&No{},
				&Available{},
				&DownloadProgress{},
				&DownloadComplete{},
				&InstallComplete{},
				&InstallRestart{},
				&Fail{},
			},
		}, ctx.Self())

		state.loop(ctx)

		log.Info("Auto update actor started")
	}
}

func (state *autoUpdateActor) checking(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *No, Fail:
		state.loop(ctx)
	case *Available:
		log.Info("Available")
		state.listener.Tell(msg)

		state.loop(ctx)

		// ctx.SetBehavior(state.downloading)

		// state.updater.Request(&Download{}, ctx.Self())
	}
}

func (state *autoUpdateActor) downloading(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Fail:
		state.loop(ctx)
	case *DownloadComplete:
		state.listener.Tell(msg)

		ctx.SetBehavior(state.installing)

		state.updater.Request(&Install{}, ctx.Self())
	}
}

func (state *autoUpdateActor) installing(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Fail:
		state.loop(ctx)
	case *InstallComplete:
		state.listener.Tell(msg)

		state.loop(ctx)
	case *InstallRestart:
	}
}

func (state *autoUpdateActor) loop(ctx actor.Context) {
	conf := config.UpdateServer()

	ctx.SetBehavior(state.checking)

	go func() {
		time.Sleep(conf.CheckInterval * time.Millisecond)

		state.updater.Request(&Check{}, ctx.Self())
	}()
}
