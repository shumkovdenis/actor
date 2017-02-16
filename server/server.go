package server

import (
	"crypto/tls"

	"go.uber.org/zap"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
	"github.com/shumkovdenis/club/manifest"
)

var log = logger.Get()

// var regPlugin *registryPlugin

// func init() {
// 	regPlugin = newRegistryPlugin()
// }

// func RegistryPlugin() *registryPlugin {
// 	return regPlugin
// }

func Start() error {
	log.Info("Start server",
		zap.String("version", manifest.Version()),
		zap.Int("port", config.Server().Port),
	)

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	var props *actor.Props

	// props = actor.FromProducer(newRoomManagerActor)
	// actor.SpawnNamed(props, "rooms")

	props = actor.FromProducer(newSessionManagerActor)
	actor.SpawnNamed(props, "sessions")

	props = actor.FromProducer(newRatesActor)
	actor.SpawnNamed(props, "rates")

	props = actor.FromProducer(newServerActor)
	actor.SpawnNamed(props, "server")

	console.ReadLine()

	return nil
}
