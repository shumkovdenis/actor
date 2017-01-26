package server

import (
	"crypto/tls"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
	"github.com/shumkovdenis/club/manifest"
	"github.com/uber-go/zap"
)

var log = logger.Get()

var regPlugin *registryPlugin

func init() {
	regPlugin = newRegistryPlugin()
}

func RegistryPlugin() *registryPlugin {
	return regPlugin
}

func Start() error {
	log.Info("Start server",
		zap.String("version", manifest.Version()),
		zap.Int("port", config.Server().Port),
	)

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	props := actor.FromProducer(newServerActor).
		WithMiddleware(plugin.Use(RegistryPlugin()))
	actor.SpawnNamed(props, "server")

	console.ReadLine()

	return nil
}
