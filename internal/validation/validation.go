package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

var v = validator.New(validator.WithRequiredStructEnabled())

func Validate[T any](input T) Errors {
	var output Errors

	rv := reflect.ValueOf(input)
	if !rv.IsValid() {
		output.Add("", "invalid input: <nil>")
		return output
	}

	// *structも引数として許可
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			output.Add("", "invalid input: nil pointer")
			return output
		}
		if rv.Elem().Kind() != reflect.Struct {
			output.Add("", "Validate expects a struct or pointer to struct")
			return output
		}
	} else if rv.Kind() != reflect.Struct {
		output.Add("", "Validate expects a struct or pointer to struct")
		return output
	}

	if err := v.Struct(input); err != nil {
		// validatorのエラーを処理
		if inv, ok := err.(*validator.InvalidValidationError); ok {
			output.Add("", inv.Error())
			return output
		}
		if verrs, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range verrs {
				output.Add(fe.StructField(), fe.Error())
			}
		} else {
			output.Add("", err.Error())
		}
	}

	if output.Has() {
		return output
	}
	return nil
}
