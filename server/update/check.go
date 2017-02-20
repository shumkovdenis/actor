package update

import (
	"net/http"

	"github.com/cavaliercoder/grab"
	"github.com/shumkovdenis/club/config"
	"go.uber.org/zap"
)

func check() Message {
	conf := config.UpdateServer()

	url := conf.PropsURL()
	path := conf.PropsPath()

	log.Info("Check update",
		zap.String("url", url),
		zap.String("path", path),
	)

	req, err := grab.NewRequest(url)
	if err != nil {
		log.Error("check update failed",
			zap.Error(err),
		)
		return &CheckFailed{}
	}

	req.Filename = path
	req.CreateMissing = true

	resp, err := grab.DefaultClient.Do(req)
	if err != nil {
		log.Error("check update failed",
			zap.Error(err),
		)
		return &CheckFailed{}
	}

	code := resp.HTTPResponse.StatusCode

	switch code {
	case http.StatusOK:
		log.Info("Update available")
		return &Available{}
	case http.StatusNotFound:
		log.Info("Update no")
		return &No{}
	default:
		log.Error("check update failed",
			zap.Int("code", code),
		)
		return &CheckFailed{}
	}
}
