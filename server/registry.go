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
}

type Registry interface {
	Register(act Act)
	Unregister(act Act)
	ToMessage(cmd *command) (interface{}, error)
	FromMessage(msg interface{}) (*event, error)
}

type Conv interface {
	SetRegistry(Registry)
}

type registry struct {
	records *hashset.Set
	acts    *treemap.Map
	cmds    *treemap.Map
}

func newRegistry() Registry {
	return &registry{
		records: hashset.New(),
		acts:    treemap.NewWithStringComparator(),
		cmds:    treemap.NewWithStringComparator(),
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

func (r *registry) ToMessage(cmd *command) (interface{}, error) {
	sample, ok := r.cmds.Get(cmd.Type)
	if !ok {
		return nil, errors.New("Command not found")
	}

	val := reflect.ValueOf(sample)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	msg := reflect.New(val.Type()).Interface().(Command)

	return msg, nil
}

func (r *registry) FromMessage(msg interface{}) (*event, error) {
	m, ok := msg.(Event)
	if !ok {
		return nil, errors.New("Message must implement Event")
	}

	evt := &event{
		Type: m.Event(),
		Data: m,
	}

	return evt, nil
}
