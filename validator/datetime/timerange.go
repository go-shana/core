package validator

import (
	"time"

	"github.com/go-shana/core/validator"
)

// TimeRange represents a time range.
type TimeRange struct {
	Min, Max time.Time

	MinInclusive bool // Whether the minimum value is included in the range.
	MaxInclusive bool // Whether the maximum value is included in the range.
}

var _ validator.Ranger[time.Time] = &TimeRange{}

// InRange returns true if the v is within the range.
func (tr *TimeRange) InRange(v time.Time) bool {
	if v.Before(tr.Min) {
		return false
	}

	if v.After(tr.Max) {
		return false
	}

	if !tr.MinInclusive && !v.After(tr.Min) {
		return false
	}

	if !tr.MaxInclusive && !v.Before(tr.Max) {
		return false
	}

	return true
}

// NewTimeRange returns a new half-open-interval TimeRange.
func NewTimeRange(min, max time.Time) *TimeRange {
	return &TimeRange{
		Min: min,
		Max: max,

		MinInclusive: true,
	}
}
