package server

import "github.com/AsynkronIT/protoactor-go/actor"

type registryPlugin struct {
	reg Registry
}

func newRegistryPlugin() *registryPlugin {
	return &registryPlugin{}
}

func (p *registryPlugin) OnStart(ctx actor.Context) {
	switch act := ctx.Actor().(type) {
	case Server:
		p.reg = act.Registry()

		log.Info("Registry plugin: server start")
	case Act:
		p.reg.Register(act)

		log.Debug("Registry plugin: register Act")
	}
}

func (p *registryPlugin) OnOtherMessage(ctx actor.Context, msg interface{}) {
	switch msg.(type) {
	case *actor.Stopped:
		if act, ok := ctx.Actor().(Act); ok {
			p.reg.Unregister(act)
		}
	}
}
