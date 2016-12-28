package update

import "github.com/AsynkronIT/gam/actor"

// Check -> command.app.update.check
type Check struct {
}

// No -> event.app.update.no
type No struct {
}

// Download -> event.app.update.download
type Download struct {
}

// Ready -> event.app.update.ready
type Ready struct {
}

// Install -> command.app.update.install
type Install struct {
}

// Restart -> event.app.update.restart
type Restart struct {
}

// Fail -> event.app.update.fail
type Fail struct {
	Message string `json:"message"`
}

type updateActor struct {
}

func NewActor() actor.Actor {
	return &updateActor{}
}

func (state *updateActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		ctx.Become(state.started)
	}
}

func (state *updateActor) started(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *Check:
	}
}

func (state *updateActor) check(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *Check:
	}
}
