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
	toMessage(cmd *Command) (interface{}, error)
	fromMessage(msg interface{}) (*Event, error)
}

type registry struct {
	recs *hashset.Set
}

func NewRegistry() Registry {
	return &registry{
		recs: hashset.New(),
	}
}

func (r *registry) AddRecord(rec Record) {
	r.recs.Add(rec)
}

func (r *registry) RemoveRecord(rec Record) {
	r.recs.Remove(rec)
}

func (r *registry) toMessage(cmd *Command) (interface{}, error) {
	var msg interface{}

	for _, rec := range r.recs.Values() {
		msg = rec.(Record).Command(cmd.Type)

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

func (r *registry) fromMessage(msg interface{}) (*Event, error) {
	var typ string

	for _, rec := range r.recs.Values() {
		typ = rec.(Record).Event(msg)

		if len(typ) > 0 {
			break
		}
	}

	if len(typ) > 0 {
		return nil, errors.New("Event not found")
	}

	evt := &Event{
		Type: typ,
		Data: msg,
	}

	return evt, nil
}
