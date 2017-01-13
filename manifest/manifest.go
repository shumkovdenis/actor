package manifest

import (
	"github.com/shumkovdenis/club/utils"
	"github.com/spf13/viper"
)

const (
	file = "manifest.json"
)

var (
	v *viper.Viper
	m *manifest
)

func init() {
	v = viper.New()
	m = new()
}

type manifest struct {
	Version string `mapstructure:"version" validate:"required"`
}

func new() *manifest {
	return &manifest{}
}

func Version() string {
	return m.Version
}

func Viper() *viper.Viper {
	return v
}

func Read() error {
	v.SetConfigType("json")
	v.SetConfigFile(file)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(m); err != nil {
		return err
	}

	if err := utils.Validate(m); err != nil {
		return err
	}

	return nil
}
