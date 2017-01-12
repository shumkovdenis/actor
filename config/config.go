package config

import (
	"fmt"
	"time"

	"github.com/shumkovdenis/club/manifest"
)

var c *config

func init() {
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

func (c *updateServer) CheckURL() string {
	return fmt.Sprintf("%s/%s/props.json", c.URL, manifest.Version())
}

func (c *updateServer) DownloadURL() string {
	return fmt.Sprintf("%s/%s/data.zip", c.URL, manifest.Version())
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

func Get() *config {
	return c
}
