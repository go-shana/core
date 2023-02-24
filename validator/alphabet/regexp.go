package alphabet

import (
	"regexp"

	"github.com/go-shana/core/validator"
)

type regexpRange struct {
	re *regexp.Regexp
}

var _ validator.Ranger[string] = regexpRange{}

func (rr regexpRange) InRange(v string) bool {
	return rr.re.MatchString(string(v))
}

// NewRegexpRange returns a new Ranger which returns true if a string matches the regular expression expr.
func NewRegexpRange(expr string) validator.Ranger[string] {
	re := regexp.MustCompile(expr)
	return regexpRange{
		re: re,
	}
}
