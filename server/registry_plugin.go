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
	case Record:
		p.reg.AddRecord(act)
	}
}

func (p *registryPlugin) OnOtherMessage(ctx actor.Context, msg interface{}) {
	switch msg.(type) {
	case *actor.Stopped:
		if rec, ok := ctx.Actor().(Record); ok {
			p.reg.RemoveRecord(rec)
		}
	}
}
