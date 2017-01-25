package plugins

import "github.com/AsynkronIT/protoactor-go/actor"

type group struct {
	set *actor.PIDSet
}

func (g group) Tell(msg interface{}) {
	g.set.ForEach(func(i int, pid actor.PID) {
		pid.Tell(msg)
	})
}
