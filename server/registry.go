package server

import (
	"errors"

	"reflect"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/hashset"
)

type Act interface {
	Name() string
	Commands() []Command
	Events() []Event
}

type Record interface {
	Command(typ string) interface{}
	Event(msg interface{}) string
}

type Registry interface {
	Register(act Act)
	Unregister(act Act)
	AddRecord(rec Record)
	RemoveRecord(rec Record)
	ToMessage(cmd *command) (Command, error)
	FromMessage(msg Event) (*event, error)
}

type registry struct {
	records *hashset.Set
	acts    *treemap.Map
	cmds    *treemap.Map
}

func newRegistry() Registry {
	return &registry{
		records: hashset.New(),
	}
}

func (r *registry) Register(act Act) {
	if c, ok := r.acts.Get(act.Name()); ok {
		r.acts.Put(act.Name(), c.(int)+1)

		return
	}

	r.acts.Put(act.Name(), 1)

	for _, cmd := range act.Commands() {
		r.cmds.Put(cmd.Command(), cmd)
	}
}

func (r *registry) Unregister(act Act) {
	c, ok := r.acts.Get(act.Name())
	if !ok {
		return
	}

	if c.(int)-1 > 0 {
		r.acts.Put(act.Name(), c.(int)-1)

		return
	}

	r.acts.Remove(act.Name())

	for _, cmd := range act.Commands() {
		r.cmds.Remove(cmd.Command())
	}
}

func (r *registry) AddRecord(rec Record) {
	r.records.Add(rec)
}

func (r *registry) RemoveRecord(rec Record) {
	r.records.Remove(rec)
}

func (r *registry) ToMessage(cmd *command) (Command, error) {
	sample, ok := r.cmds.Get(cmd.Type)
	if !ok {
		return nil, errors.New("Command not found")
	}

	val := reflect.ValueOf(sample)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	msg := reflect.New(val.Type()).Interface().(Command)

	// var msg interface{}

	// for _, record := range r.records.Values() {
	// 	msg = record.(Record).Command(cmd.Type)

	// 	if msg != nil {
	// 		break
	// 	}
	// }

	// if msg == nil {
	// 	return nil, errors.New("Command not found")
	// }

	// if err := mapstructure.Decode(cmd.Data, msg); err != nil {
	// 	return nil, err
	// }

	// return msg, nil
}

func (r *registry) FromMessage(msg Event) (*event, error) {
	// var typ string

	// for _, record := range r.records.Values() {
	// 	typ = record.(Record).Event(msg)

	// 	if len(typ) > 0 {
	// 		break
	// 	}
	// }

	// if len(typ) == 0 {
	// 	return nil, errors.New("Event not found")
	// }

	// evt := &event{
	// 	Type: typ,
	// 	Data: msg,
	// }

	// return evt, nil
}
