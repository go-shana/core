package validator

// Ranger represents a data which can be checked if it is in a range.
type Ranger[T any] interface {
	InRange(v T) bool
}

// Range represents a range of values.
type Range[T any] struct {
	Min, Max Comparable[T]

	MinInclusive bool // Whether the minimum value is inclusive.
	MaxInclusive bool // Whether the maximum value is inclusive.
}

var _ Ranger[int] = &Range[int]{}

// InRange returns true if v is in the range.
func (r *Range[T]) InRange(v T) bool {
	compMin := r.Min.Compare(v)

	if compMin < 0 {
		return false
	}

	compMax := r.Max.Compare(v)

	if compMax > 0 {
		return false
	}

	if compMin == 0 && !r.MinInclusive {
		return false
	}

	if compMax == 0 && !r.MaxInclusive {
		return false
	}

	return true
}

// NewRange returns a new half-open-interval Range.
func NewRange[T any](min, max Comparable[T]) *Range[T] {
	return &Range[T]{
		Min: min,
		Max: max,

		MinInclusive: true,
	}
}
