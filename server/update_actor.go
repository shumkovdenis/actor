package server

import (
	"fmt"
	"net/http"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/cavaliercoder/grab"
	"github.com/shumkovdenis/club/config"
	"github.com/uber-go/zap"
)

type updateActor struct {
}

func newUpdateActor() actor.Actor {
	return &updateActor{}
}

func (state *updateActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *UpdateCheck:
	}
}

func procCheckUpdate() interface{} {
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

		// tell(&Fail{"Update check failed"})

		return &UpdateCheckFail{}
	}

	req.Filename = path
	req.CreateMissing = true

	resp, err := grab.DefaultClient.Do(req)
	if err != nil {
		log.Error(err.Error())

		// tell(&Fail{"Update check failed"})

		return &UpdateCheckFail{}
	}

	code := resp.HTTPResponse.StatusCode

	switch code {
	case http.StatusOK:
		log.Info("Update available")

		// tell(&Available{})
		return &UpdateCheckAvailable{}
	case http.StatusNotFound:
		log.Info("Update no")

		// tell(&No{})
		return &UpdateCheckNo{}
	default:
		log.Error(fmt.Sprintf("Status code: %d", code))

		// tell(&Fail{"Update check failed"})
		return &UpdateCheckFail{}
	}
}

type UpdateCheck struct {
}

func (*UpdateCheck) Command() string {
	return "command.update.check"
}

type UpdateCheckAvailable struct {
}

func (*UpdateCheckAvailable) Event() string {
	return "event.update.check.available"
}

type UpdateCheckNo struct {
}

func (*UpdateCheckNo) Event() string {
	return "event.update.check.no"
}

type UpdateCheckFail struct {
}

func (*UpdateCheckFail) Event() string {
	return "event.update.check.fail"
}
