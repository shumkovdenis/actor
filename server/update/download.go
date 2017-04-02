package update

import (
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
	"go.uber.org/zap"
)

func download() chan Message {
	ch := make(chan Message)

	go func() {
		conf := config.UpdateServer()

		url := conf.DataURL()
		path := conf.DataPath()

		logger.L().Info("download update",
			zap.String("url", url),
			zap.String("path", path),
		)

		req, err := grab.NewRequest(url)
		if err != nil {
			logger.L().Error("download update failed",
				zap.Error(err),
			)
			ch <- &DownloadFailed{}
			close(ch)
			return
		}

		req.Filename = path
		req.CreateMissing = true

		respch := grab.DefaultClient.DoAsync(req)

		resp := <-respch

		for !resp.IsComplete() {
			ch <- &Progress{resp.Progress()}

			time.Sleep(200 * time.Millisecond)
		}

		if resp.Error != nil {
			logger.L().Error("download update failed",
				zap.Error(err),
			)
			ch <- &DownloadFailed{}
			close(ch)
			return
		}

		logger.L().Info("download update completed")

		ch <- &Complete{}
		close(ch)
	}()

	return ch
}
