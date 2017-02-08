package server

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
	uuid "github.com/satori/go.uuid"
	"github.com/uber-go/zap"
)

const (
	RoomNotFound = iota
	RoomFull
)

type RoomError struct {
	Code int
}

func (e *RoomError) Error() string {
	return fmt.Sprintf("room manager error: %d", e.Code)
}

type RoomConf struct {
	MaxPlaces  int `json:"max_places"`
	BusyPlaces int `json:"busy_places"`
}

type Room struct {
	ID   string    `json:"id"`
	Conf *RoomConf `json:"conf"`
}

type RoomManager interface {
	Create(conf *RoomConf) (*Room, error)
	Get(id string) (*Room, error)
	JoinRoom(id string) error
	LeaveRoom(id string) error
}

type roomManager struct {
	rooms *treemap.Map
}

func newRoomManager() RoomManager {
	return &roomManager{
		rooms: treemap.NewWithStringComparator(),
	}
}

func (m *roomManager) Create(conf *RoomConf) (*Room, error) {
	conf.MaxPlaces = 2
	conf.BusyPlaces = 0

	id := uuid.NewV4().String()

	room := &Room{id, conf}

	m.rooms.Put(id, room)

	return room, nil
}

func (m *roomManager) Get(id string) (*Room, error) {
	room, ok := m.rooms.Get(id)
	if !ok {
		log.Warn("room manager get: room not found", zap.String("room", id))

		return nil, &RoomError{RoomNotFound}
	}

	return room.(*Room), nil
}

func (m *roomManager) JoinRoom(id string) error {
	room, ok := m.rooms.Get(id)
	if !ok {
		log.Warn("join room: room not found", zap.String("room", id))

		return &RoomError{RoomNotFound}
	}

	conf := room.(*Room).Conf

	if conf.BusyPlaces == conf.MaxPlaces {
		log.Warn("join room: room full", zap.String("room", id))

		return &RoomError{RoomFull}
	}

	conf.BusyPlaces++

	return nil
}

func (m *roomManager) LeaveRoom(id string) error {
	room, ok := m.rooms.Get(id)
	if !ok {
		log.Warn("leave room: room not found", zap.String("room", id))

		return &RoomError{RoomNotFound}
	}

	conf := room.(*Room).Conf

	conf.BusyPlaces--

	return nil
}
