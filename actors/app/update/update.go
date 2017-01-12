package update

import (
	"fmt"
	"time"

	"net/http"

	"github.com/AsynkronIT/gam/actor"
	"github.com/cavaliercoder/grab"
	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/actors"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
	"github.com/uber-go/zap"
)

var log = logger.Get()

type updateActor struct {
	// listener *actor.PID
}

func New() actor.Actor {
	return &updateActor{}
}

func (state *updateActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		ctx.Become(state.started)

		log.Info("Update started")
	}
}

func (state *updateActor) started(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *Check:
		log.Info("Check update")

		ok, err := check()
		if err != nil {
			log.Error(err.Error())

			ctx.Respond(&Fail{err.Error()})

			return
		}

		if ok {
			log.Info("Update available")

			ctx.Respond(&Available{})
		} else {
			log.Info("Update no")

			ctx.Respond(&No{})
		}
	case *Download:
		actors.Process(download, ctx.Respond)
	case *Install:
	}
}

func check() (bool, error) {
	resp, err := resty.R().Get(config.UpdateServer().CheckURL())
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

func download(tell actors.Tell) {
	url := config.UpdateServer().DownloadURL()

	log.Info("Download update", zap.String("url", url))

	respch, err := grab.GetAsync(".", url)
	if err != nil {
		log.Error(err.Error())

		tell(&Fail{err.Error()})

		return
	}

	resp := <-respch

	for !resp.IsComplete() {
		tell(&DownloadProgress{resp.Progress()})

		time.Sleep(200 * time.Millisecond)
	}

	if resp.Error != nil {
		log.Error(resp.Error.Error())

		tell(&Fail{resp.Error.Error()})

		return
	}

	log.Info("Download update complete")

	tell(&DownloadComplete{})
}

func install() {

}
