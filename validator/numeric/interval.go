package numeric

import (
	"github.com/go-shana/core/validator"
)

// Numeric represents a numeric type.
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type intervalRange[T Numeric] struct {
	min, max     T
	minInclusive bool
	maxInclusive bool
}

var _ validator.Ranger[int] = &intervalRange[int]{}

func (r *intervalRange[T]) InRange(v T) bool {
	if v > r.min && v < r.max {
		return true
	}

	if v == r.min && r.minInclusive {
		return true
	}

	if v == r.max && r.maxInclusive {
		return true
	}

	return false
}

// NewOpenInterval creates a new open-interval Ranger.
func NewOpenInterval[T Numeric](min, max T) validator.Ranger[T] {
	return &intervalRange[T]{
		min: min,
		max: max,
	}
}

// NewHalfOpenInterval creates a new half-open-interval Ranger.
func NewHalfOpenInterval[T Numeric](min, max T) validator.Ranger[T] {
	return &intervalRange[T]{
		min:          min,
		max:          max,
		minInclusive: true,
	}
}

// NewClosedInterval creates a new closed-interval Ranger.
func NewClosedInterval[T Numeric](min, max T) validator.Ranger[T] {
	return &intervalRange[T]{
		min:          min,
		max:          max,
		minInclusive: true,
		maxInclusive: true,
	}
}
