package numeric

import "github.com/go-shana/core/errors"

// InRange returns true if the value is in the range [min, max).
func InRange[T Numeric](v, min, max T) {
	if v >= min && v < max {
		return
	}

	errors.Throw(errNumericValidator)
}

// InRangInclusive returns true if the value is in the range [min, max].
func InRangInclusive[T Numeric](v, min, maxInclusive T) {
	if v >= min && v <= maxInclusive {
		return
	}

	errors.Throw(errNumericValidator)
}

// LessThan returns true if the v < max.
func LessThan[T Numeric](v, max T) {
	if v < max {
		return
	}

	errors.Throw(errNumericValidator)
}

// LessEqualThan returns true if the v <= maxInclusive.
func LessEqualThan[T Numeric](v, maxInclusive T) {
	if v <= maxInclusive {
		return
	}

	errors.Throw(errNumericValidator)
}

// GreaterThan returns true if the v > min.
func GreaterThan[T Numeric](v, min T) {
	if v > min {
		return
	}

	errors.Throw(errNumericValidator)
}

// GreaterEqualThan returns true if the v >= minInclusive.
func GreaterEqualThan[T Numeric](v, minInclusive T) {
	if v >= minInclusive {
		return
	}

	errors.Throw(errNumericValidator)
}
