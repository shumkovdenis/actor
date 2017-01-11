package manifest

import (
	"errors"
	"reflect"
	"time"

	validator "gopkg.in/go-playground/validator.v9"

	"fmt"
	"strings"
)

type accountAPI struct {
	URL  string `mapstructure:"url" validate:"required,url"`
	Type string `mapstructure:"type" validate:"eq=ALLIN|eq=BINOPT"`
}

type ratesAPI struct {
	URL         string        `mapstructure:"url" validate:"required,url"`
	GetInterval time.Duration `mapstructure:"get_interval" validate:"min=1000"`
}

type updateServer struct {
	URL           string        `mapstructure:"url" validate:"required,url"`
	CheckInterval time.Duration `mapstructure:"check_interval" validate:"min=5000"`
}

type config struct {
	AccountAPI   *accountAPI   `mapstructure:"account_api"`
	RatesAPI     *ratesAPI     `mapstructure:"rates_api"`
	UpdateServer *updateServer `mapstructure:"update_server"`
}

func newConfig() *config {
	return &config{
		AccountAPI: &accountAPI{},
		RatesAPI: &ratesAPI{
			GetInterval: 5000,
		},
		UpdateServer: &updateServer{
			CheckInterval: 5000,
		},
	}
}

func (c *config) validate() error {
	validate := validator.New()

	if err := validate.Struct(c); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("fail validate config")
		}

		errs := err.(validator.ValidationErrors)
		if len(errs) > 0 {
			m := walk(reflect.TypeOf(c).Elem())
			err := errs[0]
			name := strings.TrimPrefix(err.Namespace(), "config.")
			name = m[name]
			tag := err.Tag()

			return fmt.Errorf("config '%s' failed on the '%s' tag", name, tag)
		}
	}

	c.RatesAPI.GetInterval *= time.Millisecond
	c.UpdateServer.CheckInterval *= time.Millisecond

	return nil
}

func walk(t reflect.Type) map[string]string {
	m := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		tag := field.Tag.Get("mapstructure")

		if field.Type.Kind() == reflect.Ptr {
			sm := walk(field.Type.Elem())

			for k, v := range sm {
				m[name+"."+k] = tag + "." + v
			}
		} else {
			m[name] = tag
		}
	}

	return m
}
