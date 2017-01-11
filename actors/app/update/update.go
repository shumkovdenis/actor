package update

import (
	"fmt"
	"time"

	"net/http"

	"github.com/AsynkronIT/gam/actor"
	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/actors/group"
	"github.com/shumkovdenis/club/manifest"
)

// Check -> command.app.update.check
type Check struct {
}

// No -> event.app.update.no
type No struct {
}

// Available -> event.app.update.available
type Available struct {
}

// Download -> event.app.update.download
type Download struct {
	Progress float64 `json:"progress"`
}

// Ready -> event.app.update.ready
type Ready struct {
}

// Install -> command.app.update.install
type Install struct {
}

// Restart -> event.app.update.restart
type Restart struct {
}

// Fail -> event.app.update.fail
type Fail struct {
	Message string `json:"message"`
}

type updateActor struct {
	checkURL      string
	checkInterval time.Duration
	downloadURL   string
	listener      *actor.PID
}

func NewActor(listener *actor.PID) actor.Actor {
	m := manifest.Get()
	url := m.Config.UpdateServer.URL + "/" + m.Version + "/"
	return &updateActor{
		checkURL:      url + "props.json",
		checkInterval: m.Config.UpdateServer.CheckInterval,
		downloadURL:   url + "data.zip",
		listener:      listener,
	}
}

func (state *updateActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		state.listener.Request(&group.Use{
			Producer: ctx.Self(),
			Types: []interface{}{
				&No{},
				&Available{},
				&Download{},
				&Ready{},
				&Install{},
				&Restart{},
				&Fail{},
			},
		}, ctx.Self())

		go state.checkUpdateLoop()
	}
}

func (state *updateActor) checkUpdateLoop() {
	ticker := time.Tick(state.checkInterval)
	for _ = range ticker {
		ok, err := state.checkUpdate()
		if err != nil {
			state.listener.Tell(&Fail{err.Error()})
			continue
		}
		if ok {
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
			state.listener.Tell(&No{})
		}
	}
}

func (state *updateActor) checkUpdate() (bool, error) {
	resp, err := resty.R().
		Get(state.checkURL)
	if err != nil {
		return false, fmt.Errorf("Request fail: %s", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("Status code: %d", resp.StatusCode())
	}
}
