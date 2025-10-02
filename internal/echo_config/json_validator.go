package echo_config

import "github.com/go-playground/validator/v10"

type JSONValidator struct {
	Validator *validator.Validate
}

func (jv *JSONValidator) Validate(i interface{}) error {
	if err := jv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}
