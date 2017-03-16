package update

import (
	"os"

	"github.com/Jeffail/gabs"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/manifest"
	"github.com/shumkovdenis/club/packer"
	"go.uber.org/zap"
)

func install() Message {
	conf := config.UpdateServer()

	log.Info("install update")

	json, err := gabs.ParseJSONFile(conf.PropsPath())
	if err != nil {
		log.Error("install update failed",
			zap.Error(err),
		)
		return &InstallFailed{}
	}

	restart := json.Path("restart").Data().(bool)

	if restart {
		if err := packer.Unpack(conf.DataPath(), conf.DataDir()); err != nil {
			log.Error("install update failed",
				zap.Error(err),
			)
			return &InstallFailed{}
		}

		return &Wait{}
	}

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

	if err := os.RemoveAll(conf.UpdatePath()); err != nil {
		log.Error("remove update path failed",
			zap.Error(err),
		)
	}

	log.Info("install update completed",
		zap.Bool("restart", restart),
		zap.String("version", manifest.Version()),
	)

	return &Ready{}
}
