// Package alphabet provides a set of string validators.
package alphabet

import (
	"strings"

	"github.com/go-shana/core/validator"
)

// StringRange is a ranger for string.
type StringRange struct {
	Min, Max string

	MinInclusive bool
	MaxInclusive bool
}

var _ validator.Ranger[string] = &StringRange{}

// InRange returns true if the v is within the range.
func (sr *StringRange) InRange(v string) bool {
	str := string(v)
	minStr := string(sr.Min)
	compMin := strings.Compare(str, minStr)

	if compMin < 0 {
		return false
	}

	if !sr.MinInclusive && compMin == 0 {
		return false
	}

	maxStr := string(sr.Max)
	compMax := strings.Compare(str, maxStr)

	if compMax > 0 {
		return false
	}

	if !sr.MaxInclusive && compMax == 0 {
		return false
	}

	return true
}

// NewStringRange returns a new half-open-interval StringRange.
func NewStringRange(min, max string) *StringRange {
	return &StringRange{
		Min: min,
		Max: max,

		MinInclusive: true,
	}
}
