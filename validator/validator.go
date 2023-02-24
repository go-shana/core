// Validator is a package to provide lots of useful validators to validate API request data.
package validator

import (
	"context"
	"reflect"

	"github.com/go-shana/core/errors"
)

// Validator represents a data which can work with Shana validator.
type Validator interface {
	Validate(ctx context.Context)
}

var typeOfValidator = reflect.TypeOf((*Validator)(nil)).Elem()
var errValidate = errors.New("fail to validate data")

// Validate validates data and its struct fields in depth-first order.
func Validate(ctx context.Context, data any) (err error) {
	if data == nil {
		err = errValidate
		return
	}

	val := reflect.ValueOf(data)
	err = validateValue(ctx, val)
	return
}

// validateValue validates val and its struct fields in depth-first order.
func validateValue(ctx context.Context, val reflect.Value) (err error) {
	defer errors.Handle(&err)

	elem := val

	for elem.IsValid() && elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}

	if elem.Kind() == reflect.Struct {
		t := elem.Type()
		num := elem.NumField()

		// Depth-first validation.
		for i := 0; i < num; i++ {
			stField := t.Field(i)

			if !stField.IsExported() {
				continue
			}

			field := elem.Field(i)

			if !field.IsValid() {
				continue
			}

			errors.Check(validateValue(ctx, field))
		}
	}

	for val.IsValid() {
		if val.Type().Implements(typeOfValidator) {
			val.Interface().(Validator).Validate(ctx)
			break
		}

		if val.Kind() != reflect.Pointer {
			break
		}

		val = val.Elem()
	}

	return
}
