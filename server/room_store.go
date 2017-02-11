package server

import (
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/juju/errors"
)

type RoomStore interface {
	GetByID(id string) (*Room, error)
	GetAll() ([]*Room, error)
	Add(room *Room) error
	Update(room *Room) error
}

type roomStore struct {
	rooms *treemap.Map
}

func newRoomStore() RoomStore {
	return &roomStore{
		rooms: treemap.NewWithStringComparator(),
	}
}

func (s *roomStore) GetByID(id string) (*Room, error) {
	room, ok := s.rooms.Get(id)
	if !ok {
		return nil, newRoomNotFound()
	}
	return room.(*Room), nil
}

func (s *roomStore) GetAll() ([]*Room, error) {
	rooms := make([]*Room, s.rooms.Size())
	for i, room := range s.rooms.Values() {
		rooms[i] = room.(*Room)
	}
	return rooms, nil
}

func (s *roomStore) Add(room *Room) error {
	s.rooms.Put(room.ID, room)
	return nil
}

func (s *roomStore) Update(room *Room) error {
	s.rooms.Put(room.ID, room)
	return nil
}

type roomNotFound struct {
	errors.Err
}

func newRoomNotFound() error {
	err := &roomNotFound{errors.NewErr("room not found")}
	err.SetLocation(1)
	return err
}
