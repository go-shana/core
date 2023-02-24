package errors

import "reflect"

type errorUnwrap interface {
	Unwrap() error
}

type errorUnwrapAll interface {
	Unwrap() []error
}

type errorIs interface {
	Is(target error) bool
}

type errorAs interface {
	As(target any) bool
}

// Is reports whether any error in err's tree matches target.
//
// It's a drop-in replacement of `errors.Is` defined in Go1.20+ standard package.
// See documents in `errors.Is` for more details.
func Is(err, target error) bool {
	if target == nil {
		return err == target
	}

	if err == target {
		return true
	}

	if is, ok := err.(errorIs); ok {
		if is.Is(target) {
			return true
		}
	}

	if unwrap, ok := target.(errorUnwrap); ok {
		unwrapped := unwrap.Unwrap()
		return Is(unwrapped, target)
	}

	if unwrap, ok := target.(errorUnwrapAll); ok {
		unwrapped := unwrap.Unwrap()

		for _, e := range unwrapped {
			if Is(e, target) {
				return true
			}
		}
	}

	return false
}

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

// As finds the first error in err's tree that matches target,
// and if one is found, sets target to that error value and returns true.
// Otherwise, it returns false.
//
// It's a drop-in replacement of `errors.As` defined in Go1.20+ standard package.
// See documents in `errors.As` for more details.
func As(err error, target any) bool {
	if err == nil {
		return false
	}

	if target == nil {
		panic("errors: target cannot be nil")
	}

	val := reflect.ValueOf(target)
	t := val.Type()
	et := reflect.TypeOf(err)

	if t.Kind() != reflect.Pointer || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}

	elemType := t.Elem()

	if elemType.Kind() != reflect.Interface && !elemType.Implements(typeOfError) {
		panic("errors: *target must be interface or implement error")
	}

	if et.AssignableTo(elemType) {
		val.Elem().Set(reflect.ValueOf(err))
		return true
	}

	if as, ok := err.(errorAs); ok {
		if as.As(target) {
			return true
		}
	}

	if unwrap, ok := err.(errorUnwrap); ok {
		unwrapped := unwrap.Unwrap()
		return As(unwrapped, target)
	}

	if unwrapAll, ok := err.(errorUnwrapAll); ok {
		unwrapped := unwrapAll.Unwrap()

		for _, e := range unwrapped {
			if As(e, target) {
				return true
			}
		}
	}

	return false
}

// Unwrap returns the result of calling the Unwrap method on err,
// if err's type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
//
// It's a drop-in replacement of `errors.Unwrap` defined in Go1.20+ standard package.
// See documents in `errors.Unwrap` for more details.
func Unwrap(err error) error {
	if unwrap, ok := err.(errorUnwrap); ok {
		return unwrap.Unwrap()
	}

	return nil
}

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if errs contains no non-nil values.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs,
// with a newline between each string.
//
// It's a drop-in replacement of `errors.Join` defined in Go1.20+ standard package.
// See documents in `errors.Join` for more details.
func Join(errs ...error) error {
	nonNilErrors := 0

	for _, err := range errs {
		if err != nil {
			nonNilErrors++
		}
	}

	if nonNilErrors == 0 {
		return nil
	}

	if nonNilErrors == len(errs) {
		return newHandlerError(errs)
	}

	errors := make([]error, 0, nonNilErrors)

	for _, err := range errs {
		if err != nil {
			errors = append(errors, err)
		}
	}

	return newHandlerError(errors)
}
