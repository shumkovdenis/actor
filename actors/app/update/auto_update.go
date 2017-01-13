package update

import (
	"time"

	"github.com/AsynkronIT/gam/actor"
	"github.com/shumkovdenis/club/actors/group"
	"github.com/shumkovdenis/club/config"
)

type autoUpdateActor struct {
	updater  *actor.PID
	listener *actor.PID
	ticker   *time.Ticker
}

func newAutoUpdate(updater *actor.PID) actor.Actor {
	return &autoUpdateActor{
		updater: updater,
	}
}

func (state *autoUpdateActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		props := actor.FromProducer(group.New)
		state.listener = ctx.SpawnNamed(props, "auto")

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

		go state.loop(ctx.Self())

		log.Info("Auto update actor started")
	case *Available:
		state.listener.Tell(msg)

		// state.updater.Request(&Download{}, ctx.Self())
	case *DownloadComplete:
		state.listener.Tell(msg)

		state.updater.Request(&Install{}, ctx.Self())
	case
		*No,
		*DownloadProgress,
		*InstallComplete,
		*InstallRestart:
		state.listener.Tell(msg)
	}
}

func (state *autoUpdateActor) loop(self *actor.PID) {
	conf := config.UpdateServer()

	ticker := time.NewTicker(conf.CheckInterval * time.Millisecond)
	for _ = range ticker.C {
		state.updater.Request(&Check{}, self)
	}
}
