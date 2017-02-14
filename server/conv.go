package server

import (
	"errors"
	"reflect"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/mitchellh/mapstructure"
)

var commands = treemap.NewWithStringComparator()

func init() {
	cmds := []Command{
		&Subscribe{},
		&Unsubscribe{},
		&Login{},
		&AccountAuth{},
		&AccountBalance{},
		&AccountSession{},
		&AccountWithdraw{},
	}

	for _, cmd := range cmds {
		commands.Put(cmd.Command(), cmd)
	}
}

type Conv interface {
	ToMessage(cmd *command) (interface{}, error)
	FromMessage(evt interface{}) (*event, error)
}

type conv struct {
}

func newConv() Conv {
	return &conv{}
}

func (c *conv) ToMessage(cmd *command) (interface{}, error) {
	sample, ok := commands.Get(cmd.Type)
	if !ok {
		return nil, errors.New("command not found")
	}

	val := reflect.ValueOf(sample)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	msg := reflect.New(val.Type()).Interface()

	if err := mapstructure.Decode(cmd.Data, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (c *conv) FromMessage(msg interface{}) (*event, error) {
	m, ok := msg.(Event)
	if !ok {
		return nil, errors.New("message must implement 'Event'")
	}

	evt := &event{
		Type: m.Event(),
		Data: m,
	}

	return evt, nil
}
