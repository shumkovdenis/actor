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
	"github.com/shumkovdenis/club/server/jackpots/list"
	"github.com/shumkovdenis/club/server/jackpots/tops"
	"github.com/shumkovdenis/club/server/rates"
	"github.com/shumkovdenis/club/server/update"
)

var log = logger.Get()

func Start() error {
	log.Info("start server",
		zap.String("version", manifest.Version()),
		zap.Int("port", config.Server().Port),
	)

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	var props *actor.Props

	props = actor.FromProducer(update.NewActor)
	actor.SpawnNamed(props, "update")

	props = actor.FromProducer(newRoomManagerActor)
	actor.SpawnNamed(props, "rooms")

	props = actor.FromProducer(newSessionManagerActor)
	actor.SpawnNamed(props, "sessions")

	props = actor.FromProducer(tops.NewActor)
	actor.SpawnNamed(props, "jackpots/tops")

	props = actor.FromProducer(list.NewActor)
	actor.SpawnNamed(props, "jackpots/list")

	props = actor.FromProducer(rates.NewActor)
	actor.SpawnNamed(props, "rates")

	props = actor.FromProducer(newServerActor)
	actor.SpawnNamed(props, "server")

	console.ReadLine()

	return nil
}
