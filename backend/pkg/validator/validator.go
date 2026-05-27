package validator

import "github.com/go-playground/validator/v10"

var instance = validator.New()

func Validate(v interface{}) error {
	return instance.Struct(v)
}
