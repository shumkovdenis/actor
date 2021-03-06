package config

import (
	"os"
	"path"
	"time"

	"github.com/shumkovdenis/club/manifest"
	"github.com/shumkovdenis/club/utils"
	"github.com/spf13/viper"
)

const (
	File = "config.toml"

	appName   = "club"
	propsFile = "props.json"
	dataFile  = "data.zip"
)

var (
	v *viper.Viper
	c *config
)

func init() {
	v = viper.New()
	c = new()
}

type server struct {
	Port       int    `mapstructure:"port" validate:"required"`
	PublicPath string `mapstructure:"public_path" validate:"required"`
}

type accountAPI struct {
	URL  string `mapstructure:"url" validate:"required,url"`
	Type string `mapstructure:"type" validate:"eq=ALLIN|eq=BINOPT"`
}

type ratesAPI struct {
	URL         string        `mapstructure:"url" validate:"required,url"`
	GetInterval time.Duration `mapstructure:"get_interval" validate:"min=1000"`
}

type updateServer struct {
	URL           string        `mapstructure:"url" validate:"required,url"`
	AutoUpdate    bool          `mapstructure:"auto_update"`
	CheckInterval time.Duration `mapstructure:"check_interval" validate:"min=5000"`
}

func (c *updateServer) updateURL() string {
	return c.URL + "/" + manifest.Version() + "/"
}

func (c *updateServer) UpdatePath() string {
	return path.Join(os.TempDir(), appName+"-"+manifest.Version())
}

func (c *updateServer) PropsURL() string {
	return c.updateURL() + propsFile
}

func (c *updateServer) PropsPath() string {
	return path.Join(c.UpdatePath(), propsFile)
}

func (c *updateServer) DataURL() string {
	return c.updateURL() + dataFile
}

func (c *updateServer) DataPath() string {
	return path.Join(c.UpdatePath(), dataFile)
}

func (c *updateServer) AppPath() string {
	return "."
}

type config struct {
	Server       *server       `mapstructure:"server"`
	AccountAPI   *accountAPI   `mapstructure:"account_api"`
	RatesAPI     *ratesAPI     `mapstructure:"rates_api"`
	UpdateServer *updateServer `mapstructure:"update_server"`
}

func new() *config {
	return &config{
		Server: &server{
			Port:       8282,
			PublicPath: "public",
		},
		AccountAPI: &accountAPI{},
		RatesAPI: &ratesAPI{
			GetInterval: 5000,
		},
		UpdateServer: &updateServer{
			AutoUpdate:    false,
			CheckInterval: 5000,
		},
	}
}

func Server() *server {
	return c.Server
}

func AccountAPI() *accountAPI {
	return c.AccountAPI
}

func RatesAPI() *ratesAPI {
	return c.RatesAPI
}

func UpdateServer() *updateServer {
	return c.UpdateServer
}

func Viper() *viper.Viper {
	return v
}

func Read(file string) error {
	v.SetConfigType("toml")
	v.SetConfigFile(file)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(c); err != nil {
		return err
	}

	if err := utils.Validate(c); err != nil {
		return err
	}

	return nil
}
