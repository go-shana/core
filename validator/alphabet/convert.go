package alphabet

import "github.com/go-shana/core/validator"

type stringWrapperRange[T ~string] struct {
	ranger validator.Ranger[string]
}

var _ validator.Ranger[string] = stringWrapperRange[string]{}

func (swr stringWrapperRange[T]) InRange(v T) bool {
	return swr.ranger.InRange(string(v))
}

// As returns a new Ranger which converts the ranger type from string to the type T.
func As[T ~string](ranger validator.Ranger[string]) validator.Ranger[T] {
	return stringWrapperRange[T]{
		ranger: ranger,
	}
}
