package update

import (
	"fmt"
	"time"

	"net/http"

	"github.com/AsynkronIT/gam/actor"
	"github.com/cavaliercoder/grab"
	"github.com/go-resty/resty"
	"github.com/shumkovdenis/actor/actors/group"
	"github.com/shumkovdenis/actor/config"
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
	listener *actor.PID
}

func NewActor(listener *actor.PID) actor.Actor {
	return &updateActor{listener}
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
	updateServer := config.Conf.UpdateServer
	ticker := time.Tick(time.Duration(updateServer.UpdateInterval) * time.Millisecond)
	for _ = range ticker {
		ok, err := checkUpdate()
		if err != nil {
			state.listener.Tell(&Fail{err.Error()})
			continue
		}
		if ok {
			state.listener.Tell(&Available{})
		} else {
			state.listener.Tell(&No{})
		}
	}
}

func checkUpdate() (bool, error) {
	conf := config.Conf
	resp, err := resty.R().
		Get(conf.UpdateServer.URL + "/" + conf.Version)
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

func download() <-chan *grab.Response {
	respch, err := grab.GetAsync(".", conf.UpdateServer.URL+"/"+conf.Version+".zip")
	if err != nil {

	}
	return respch
}
