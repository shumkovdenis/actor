package manifest

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

var m *Manifest

func init() {
	m = newManifest()
}

type Manifest struct {
	Version string
	Config  *config
}

func newManifest() *Manifest {
	return &Manifest{
		Config: newConfig(),
	}
}

func Get() *Manifest {
	return m
}

func ReadConfig(path string) error {
	if len(strings.TrimSpace(path)) == 0 {
		return errors.New("must specify path to config file")
	}

	viper.SetConfigType("toml")
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(m.Config); err != nil {
		return err
	}

	if err := m.Config.validate(); err != nil {
		return err
	}

	return nil
}
