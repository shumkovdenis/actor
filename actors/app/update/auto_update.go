package update

import (
	"time"

	"github.com/AsynkronIT/gam/actor"
	"github.com/shumkovdenis/club/config"
)

type autoUpdateActor struct {
	listener *actor.PID
}

func (state *autoUpdateActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		// state.listener.Request(&group.Use{
		// 	Producer: ctx.Self(),
		// 	Types: []interface{}{
		// 		&No{},
		// 		&Available{},
		// 		&Download{},
		// 		&Ready{},
		// 		&Install{},
		// 		&Restart{},
		// 		&Fail{},
		// 	},
		// }, ctx.Self())

	}
}

func (state *autoUpdateActor) checkUpdateLoop() {
	ticker := time.Tick(config.UpdateServer().CheckInterval * time.Millisecond)
	for _ = range ticker {
		log.Info("Check update (auto)")

		ok, err := check()
		if err != nil {
			log.Error(err.Error())

			state.listener.Tell(&Fail{err.Error()})

			continue
		}
		if ok {
			log.Info("Update available (auto)")

			state.listener.Tell(&Available{})

			// respch, err := grab.GetAsync(".", state.downloadURL)
			// if err != nil {
			// 	state.listener.Tell(&Fail{err.Error()})
			// 	continue
			// }

			// resp := <-respch

			// for !resp.IsComplete() {
			// 	state.listener.Tell(&Download{resp.Progress()})
			// 	time.Sleep(200 * time.Millisecond)
			// }

			// if resp.Error != nil {
			// 	state.listener.Tell(&Fail{resp.Error.Error()})
			// 	continue
			// }

			// state.listener.Tell(&Ready{})
		} else {
			log.Info("Update no (auto)")

			state.listener.Tell(&No{})
		}
	}
}
