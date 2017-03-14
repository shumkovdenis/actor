package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/router"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/shumkovdenis/club/server/update"
)

type roomActor struct {
	sessions  *hashset.Set
	updatePID *actor.PID
}

func newRoomActor() actor.Actor {
	return &roomActor{
		sessions: hashset.New(),
	}
}

func (state *roomActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.updatePID = actor.NewLocalPID("update")
	case *JoinRoom:
		if state.sessions.Size() == 2 {
			ctx.Respond(&RoomFull{})
			return
		}

		state.sessions.Add(msg.SessionPID)

		ctx.Respond(&JoinedRoom{})
	case *LeaveRoom:
		state.sessions.Remove(msg.SessionPID)

		ctx.Respond(&LeftRoom{})
	case update.Message:
		pids := make([]*actor.PID, 0, state.sessions.Size()+1)
		pids = append(pids, ctx.Sender())
		for _, session := range state.sessions.Values() {
			pids = append(pids, session.(*actor.PID))
		}

		pid := actor.Spawn(router.NewBroadcastGroup(pids...))

		state.updatePID.Request(msg, pid)
	}
}
