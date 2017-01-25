package update

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/plugins"
)

type autoUpdateActor struct {
	brokerList *plugins.List
}

func newAutoUpdate(brokerList *plugins.List) actor.Actor {
	return &autoUpdateActor{
		brokerList: brokerList,
	}
}

func (state *autoUpdateActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		// state.listener.Request(&group.Use{
		// 	Producer: ctx.Self(),
		// Types: []interface{}{
		// 	&No{},
		// 	&Available{},
		// 	&DownloadProgress{},
		// 	&DownloadComplete{},
		// 	&InstallComplete{},
		// 	&InstallRestart{},
		// 	&Fail{},
		// },
		// }, ctx.Self())

		state.loop(ctx)

		log.Info("Auto update actor started")
	}
}

func (state *autoUpdateActor) checking(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *No, Fail:
		state.loop(ctx)
	case *Available:
		state.brokerList.Tell(msg)

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
		state.brokerList.Tell(msg)

		ctx.SetBehavior(state.installing)

		ctx.Parent().Request(&Install{}, ctx.Self())
	}
}

func (state *autoUpdateActor) installing(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Fail:
		state.loop(ctx)
	case *InstallComplete:
		state.brokerList.Tell(msg)

		state.loop(ctx)
	case *InstallRestart:
	}
}

func (state *autoUpdateActor) loop(ctx actor.Context) {
	conf := config.UpdateServer()

	ctx.SetBehavior(state.checking)

	go func() {
		time.Sleep(conf.CheckInterval * time.Millisecond)

		ctx.Parent().Request(&Check{}, ctx.Self())
	}()
}
