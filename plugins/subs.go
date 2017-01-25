package plugins

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/shumkovdenis/club/messages"
)

// Subs ...
type Subs struct {
	*group
	topic string
}

// NewSubs ...
func NewSubs(topic string) *Subs {
	return &Subs{
		group: &group{
			set: actor.NewPIDSet(),
		},
		topic: topic,
	}
}

// OnStart ...
func (s *Subs) OnStart(ctx actor.Context) {}

// OnOtherMessage ...
func (s *Subs) OnOtherMessage(ctx actor.Context, msg interface{}) {
	switch msg := ctx.Message().(type) {
	case *messages.SubscribeSuccess:
		if msg.Contains(s.topic) {
			s.set.Add(ctx.Self())

			log.Debug(fmt.Sprintf("Subs '%s': add '%s'", s.topic, ctx.Self().Id))
		}
	case *messages.UnsubscribeSuccess:
		if msg.Contains(s.topic) {
			s.set.Remove(ctx.Self())

			log.Debug(fmt.Sprintf("List '%s': remove '%s'", s.topic, ctx.Self().Id))
		}
	}
}
