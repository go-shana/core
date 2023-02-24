package validator

type enumRange[T comparable] struct {
	enums map[T]struct{}
}

var _ Ranger[int] = enumRange[int]{}

func (enum enumRange[T]) InRange(v T) (ok bool) {
	_, ok = enum.enums[v]
	return
}

// NewEnum returns a new Ranger which returns true if a value is one of the values.
func NewEnum[T comparable](values ...T) Ranger[T] {
	enums := make(map[T]struct{}, len(values))

	for _, val := range values {
		enums[val] = struct{}{}
	}

	return enumRange[T]{
		enums: enums,
	}
}
