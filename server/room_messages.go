package server

import "github.com/AsynkronIT/protoactor-go/actor"

type RoomMessage interface {
	RoomMessage()
}

type CreateRoom struct {
	RoomID string
}

func (*CreateRoom) RoomMessage() {}

type Room struct {
	ID string
}

func (*Room) RoomMessage() {}

type GetRoom struct {
	RoomID string
}

func (*GetRoom) RoomMessage() {}

type RoomNotFound struct{}

func (*RoomNotFound) RoomMessage() {}

func (*RoomNotFound) Code() string {
	return "room_not_found"
}

type JoinRoom struct {
	SessionPID *actor.PID
}

func (*JoinRoom) RoomMessage() {}

type JoinedRoom struct{}

func (*JoinedRoom) RoomMessage() {}

type RoomFull struct{}

func (*RoomFull) RoomMessage() {}

func (*RoomFull) Code() string {
	return "room_full"
}

type LeaveRoom struct {
	SessionPID *actor.PID
}

func (*LeaveRoom) RoomMessage() {}

type LeftRoom struct{}

func (*LeftRoom) RoomMessage() {}
