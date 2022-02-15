package utils

import (
	"github.com/go-playground/validator/v10"
)

type IError struct {
	Field string
	Tag   string
	Value string
}

func ValidatorErrorHandler(err error) []*IError {

	var errors []*IError

	for _, err := range err.(validator.ValidationErrors) {
		var el IError
		el.Field = err.Field()
		el.Tag = err.Tag()
		el.Value = err.Param()
		errors = append(errors, &el)
	}

	return errors
}

// func ErrorHandler(err error){
// 	var el IError
// 		el.Field = err.Error().
// 		el.Tag = err.Tag()
// 		el.Value = err.Param()
// }
