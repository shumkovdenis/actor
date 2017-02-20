package update

import (
	"os"

	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/manifest"
	"github.com/shumkovdenis/club/packer"
	"go.uber.org/zap"
)

func install() Message {
	conf := config.UpdateServer()

	log.Info("install update")

	if err := packer.Unpack(conf.DataPath(), conf.AppPath()); err != nil {
		log.Error("install update failed",
			zap.Error(err),
		)
		return &InstallFailed{}
	}

	if err := manifest.Read(); err != nil {
		log.Error("install update failed",
			zap.Error(err),
		)
		return &InstallFailed{}
	}

	if err := os.Remove(conf.UpdatePath()); err != nil {
		log.Error("install update failed",
			zap.Error(err),
		)
		return &InstallFailed{}
	}

	log.Info("install update completed",
		zap.String("version", manifest.Version()),
	)

	return &Ready{}
}
