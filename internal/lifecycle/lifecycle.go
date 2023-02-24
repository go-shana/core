package lifecycle

import (
	"context"

	"github.com/go-shana/core/errors"
)

// Func is a callback function in life cycle.
type Func func(context.Context) error

// LifeCycle is a phase in service life cycle.
type LifeCycle struct {
	Name  string
	Funcs []Func
}

// New creates a new life cycle.
func New(name string) *LifeCycle {
	return &LifeCycle{
		Name: name,
	}
}

// AddFunc adds a function to the life cycle.
func (lc *LifeCycle) AddFunc(f Func) {
	lc.Funcs = append(lc.Funcs, f)
}

// Run runs all functions in the life cycle.
func (lc *LifeCycle) Run(ctx context.Context) (err error) {
	defer errors.Handle(&err)

	for _, f := range lc.Funcs {
		errors.Check(f(ctx))
	}

	return
}

// Reset discards all callbacks.
func (lc *LifeCycle) Reset() {
	lc.Funcs = nil
}
