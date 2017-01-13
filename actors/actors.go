package actors

import (
	"github.com/AsynkronIT/gam/actor"
)

type Tell func(interface{})

type Proc func(Tell)

func Process(proc Proc, tell Tell) {
	ch := make(chan interface{})

	t := func(m interface{}) {
		ch <- m
	}

	go func() {
		proc(t)
		close(ch)
	}()

	for m := range ch {
		tell(m)
	}
}

func IsActorMessage(msg interface{}) bool {
	switch msg.(type) {
	case
		*actor.Restarting,
		*actor.Stopping,
		*actor.Stopped,
		*actor.PoisonPill,
		*actor.Started,
		*actor.Receive:
		return true
	}
	return false
}
