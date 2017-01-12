package utils

import (
	"fmt"
	"reflect"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

func Validate(s interface{}) error {
	validate := validator.New()

	if err := validate.Struct(s); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("fail validate")
		}

		errs := err.(validator.ValidationErrors)

		if len(errs) > 0 {
			err := errs[0]
			typ := reflect.TypeOf(s).Elem()
			prefix := fmt.Sprintf("%s.", typ.Name())
			name := strings.TrimPrefix(err.Namespace(), prefix)
			m := walk(typ)
			name = m[name]

			return fmt.Errorf("field validation for '%s' failed on the '%s' tag", name, err.Tag())
		}
	}

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
