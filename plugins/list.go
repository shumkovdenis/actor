package plugins

import "github.com/AsynkronIT/protoactor-go/actor"
import "fmt"

// List ...
type List struct {
	*group
	name string
}

// NewList ...
func NewList(name string) *List {
	return &List{
		group: &group{
			set: actor.NewPIDSet(),
		},
		name: name,
	}
}

// OnStart ...
func (l *List) OnStart(ctx actor.Context) {
	l.set.Add(ctx.Self())

	log.Debug(fmt.Sprintf("List '%s': add '%s'", l.name, ctx.Self().Id))
}

// OnOtherMessage ...
func (l *List) OnOtherMessage(ctx actor.Context, msg interface{}) {
	switch msg.(type) {
	case *actor.Stopped:
		l.set.Remove(ctx.Self())

		log.Debug(fmt.Sprintf("List '%s': remove '%s'", l.name, ctx.Self().Id))
	}
}
