package server

import (
	"errors"
	"reflect"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/mitchellh/mapstructure"
	"github.com/shumkovdenis/club/server/account"
	"github.com/shumkovdenis/club/server/core"
	"github.com/shumkovdenis/club/server/update"
)

var commands = treemap.NewWithStringComparator()

func init() {
	cmds := []core.Command{
		&Subscribe{},
		&Unsubscribe{},
		&Login{},
		&update.Check{},
		&update.Download{},
		&update.Install{},
		&update.Restart{},
		&account.Authorize{},
		&account.GetBalance{},
		&account.GetGameSession{},
		&account.Withdraw{},
		&account.Cashback{},
		&account.GetJackpotsTops{},
		&account.GetJackpotsList{},
	}

	for _, cmd := range cmds {
		commands.Put(cmd.Command(), cmd)
	}
}

type command struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Conv interface {
	ToMessage(cmd *command) (interface{}, error)
	FromMessage(evt interface{}) (*event, error)
}

type conv struct{}

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
	m, ok := msg.(core.Event)
	if !ok {
		return nil, errors.New("message must implement 'Event'")
	}

	var data interface{}

	if code, ok := msg.(core.Code); ok {
		data = &struct {
			Code string `json:"code"`
		}{
			Code: code.Code(),
		}
	} else {
		data = msg
	}

	evt := &event{
		Type: m.Event(),
		Data: data,
	}

	return evt, nil
}
