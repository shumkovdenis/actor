package update

import (
	"net/http"
	"os"

	"github.com/cavaliercoder/grab"
	"github.com/shumkovdenis/club/config"
	"go.uber.org/zap"
)

func check() Message {
	conf := config.UpdateServer()

	url := conf.PropsURL()
	path := conf.PropsPath()

	log.Info("check update",
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

	if code == http.StatusOK {
		log.Info("update available")
		return &Available{}
	}

	if err := os.RemoveAll(conf.UpdatePath()); err != nil {
		log.Error("remove update path failed",
			zap.Error(err),
		)
	}

	if code == http.StatusNotFound || code == http.StatusForbidden {
		log.Info("update no")
		return &No{}
	}

	log.Error("check update failed",
		zap.Int("code", code),
	)
	return &CheckFailed{}
}
