package update

import (
	"net/http"
	"os"

	"github.com/cavaliercoder/grab"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
	"go.uber.org/zap"
)

func check() Message {
	conf := config.UpdateServer()

	url := conf.PropsURL()
	path := conf.PropsPath()

	logger.L().Info("check update",
		zap.String("url", url),
		zap.String("path", path),
	)

	if err := os.RemoveAll(conf.UpdatePath()); err != nil {
		logger.L().Error("remove old update failed",
			zap.Error(err),
		)
		return &CheckFailed{}
	}

	req, err := grab.NewRequest(url)
	if err != nil {
		logger.L().Error("check update failed",
			zap.Error(err),
		)
		return &CheckFailed{}
	}

	req.Filename = path
	req.CreateMissing = true

	resp, err := grab.DefaultClient.Do(req)
	if err != nil {
		logger.L().Error("check update failed",
			zap.Error(err),
		)
		return &CheckFailed{}
	}

	code := resp.HTTPResponse.StatusCode

	if code == http.StatusOK {
		logger.L().Info("update available")
		return &Available{}
	}

	if code == http.StatusNotFound || code == http.StatusForbidden {
		logger.L().Info("update no")
		return &No{}
	}

	logger.L().Error("check update failed",
		zap.Int("code", code),
	)
	return &CheckFailed{}
}
