package server

import (
	"errors"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/mitchellh/mapstructure"
)

type Record interface {
	Command(typ string) interface{}
	Event(msg interface{}) string
}

type Registry interface {
	AddRecord(rec Record)
	RemoveRecord(rec Record)
	ToMessage(cmd *command) (interface{}, error)
	FromMessage(msg interface{}) (*event, error)
}

type registry struct {
	records *hashset.Set
}

func newRegistry() Registry {
	return &registry{
		records: hashset.New(),
	}
}

func (r *registry) AddRecord(rec Record) {
	r.records.Add(rec)
}

func (r *registry) RemoveRecord(rec Record) {
	r.records.Remove(rec)
}

func (r *registry) ToMessage(cmd *command) (interface{}, error) {
	var msg interface{}

	for _, record := range r.records.Values() {
		msg = record.(Record).Command(cmd.Type)

		if msg != nil {
			break
		}
	}

	if msg == nil {
		return nil, errors.New("Command not found")
	}

	if err := mapstructure.Decode(cmd.Data, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (r *registry) FromMessage(msg interface{}) (*event, error) {
	var typ string

	for _, record := range r.records.Values() {
		typ = record.(Record).Event(msg)

		if len(typ) > 0 {
			break
		}
	}

	if len(typ) == 0 {
		return nil, errors.New("Event not found")
	}

	evt := &event{
		Type: typ,
		Data: msg,
	}

	return evt, nil
}
