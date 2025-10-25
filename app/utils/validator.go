package utils

import (
	"github.com/go-playground/validator/v10"
)

var (
	Validator = validator.New()
)

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

func ParseValidationError(err validator.ValidationErrors) []ValidationError {
	errors := make([]ValidationError, len(err))
	for i, e := range err {
		errors[i] = ValidationError{
			Field: e.Field(),
			Tag:   e.Tag(),
		}
	}

	return errors
}
