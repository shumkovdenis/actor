package update

import (
	"fmt"
	"time"

	"net/http"

	"os"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/cavaliercoder/grab"
	"github.com/shumkovdenis/club/actors"
	"github.com/shumkovdenis/club/actors/group"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
	"github.com/shumkovdenis/club/manifest"
	"github.com/shumkovdenis/club/packer"
	"github.com/uber-go/zap"
)

var log = logger.Get()

type updateActor struct {
}

func New() actor.Actor {
	return &updateActor{}
}

func (state *updateActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		if config.UpdateServer().AutoUpdate {
			props := actor.FromProducer(group.New)
			pid, err := ctx.SpawnNamed(props, "auto")
			if err != nil {
				log.Error(err.Error())
			}

			props = actor.FromInstance(newAutoUpdate(ctx.Self(), pid))
			ctx.Spawn(props)
		}

		ctx.SetBehavior(state.started)

		log.Info("Update actor started")
	}
}

func (state *updateActor) started(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *Check:
		ctx.SetBehavior(state.check)

		actors.Process(check, ctx.Respond)

		ctx.SetBehavior(state.started)
	case *Download:
		ctx.SetBehavior(state.download)

		actors.Process(download, ctx.Respond)

		ctx.SetBehavior(state.started)
	case *Install:
		ctx.SetBehavior(state.install)

		actors.Process(install, ctx.Respond)

		ctx.SetBehavior(state.started)
	}
}

func (state *updateActor) check(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *Check:
		ctx.Respond(&Fail{"Already checking updates"})
	case
		*Download,
		*Install:
		ctx.Respond(&Fail{"Checking updates"})
	}
}

func (state *updateActor) download(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *Download:
		ctx.Respond(&Fail{"Already downloading updates"})
	case
		*Check,
		*Install:
		ctx.Respond(&Fail{"Downloading updates"})
	}
}

func (state *updateActor) install(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *Install:
		ctx.Respond(&Fail{"Already installing updates"})
	case
		*Check,
		*Download:
		ctx.Respond(&Fail{"Installing updates"})
	}
}

func check(tell actors.Tell) {
	conf := config.UpdateServer()

	url := conf.PropsURL()
	path := conf.PropsPath()

	log.Info("Check update",
		zap.String("url", url),
		zap.String("path", path),
	)

	req, err := grab.NewRequest(url)
	if err != nil {
		log.Error(err.Error())

		tell(&Fail{"Update check failed"})

		return
	}

	req.Filename = path
	req.CreateMissing = true

	resp, err := grab.DefaultClient.Do(req)
	if err != nil {
		log.Error(err.Error())

		tell(&Fail{"Update check failed"})

		return
	}

	code := resp.HTTPResponse.StatusCode

	switch code {
	case http.StatusOK:
		log.Info("Update available")

		tell(&Available{})
	case http.StatusNotFound:
		log.Info("Update no")

		tell(&No{})
	default:
		log.Error(fmt.Sprintf("Status code: %d", code))

		tell(&Fail{"Update check failed"})
	}
}

func download(tell actors.Tell) {
	conf := config.UpdateServer()

	url := conf.DataURL()
	path := conf.DataPath()

	log.Info("Download update",
		zap.String("url", url),
		zap.String("path", path),
	)

	req, err := grab.NewRequest(url)
	if err != nil {
		log.Error(err.Error())

		tell(&Fail{"Update download failed"})

		return
	}

	req.Filename = path
	req.CreateMissing = true

	respch := grab.DefaultClient.DoAsync(req)

	resp := <-respch

	for !resp.IsComplete() {
		tell(&DownloadProgress{resp.Progress()})

		time.Sleep(200 * time.Millisecond)
	}

	if resp.Error != nil {
		log.Error(resp.Error.Error())

		tell(&Fail{"Update download failed"})

		return
	}

	log.Info("Update download completed")

	tell(&DownloadComplete{})
}

func install(tell actors.Tell) {
	conf := config.UpdateServer()

	log.Info("Install update")

	if err := packer.Unpack(conf.DataPath(), conf.AppPath()); err != nil {
		log.Error(err.Error())

		tell(&Fail{"Update install failed"})

		return
	}

	if err := manifest.Read(); err != nil {
		log.Error(err.Error())

		tell(&Fail{"Update install failed"})

		return
	}

	if err := os.Remove(conf.UpdatePath()); err != nil {
		log.Error(err.Error())
	}

	log.Info("Update install completed",
		zap.String("version", manifest.Version()),
	)

	tell(&InstallComplete{})
}
