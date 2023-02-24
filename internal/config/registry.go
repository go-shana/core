package config

import (
	"context"

	"github.com/go-shana/core/data"
	"github.com/go-shana/core/errors"
)

const TagName = "shana"

type registry struct {
	entries []registryEntry
}

type registryEntry struct {
	Query    string
	Value    any
	DoneFunc DoneFunc
}

type DoneFunc func(context.Context) error

var defaultRegistry = registry{}

// DefaultRegistry returns the default registry.
func DefaultRegistry() *registry {
	return &defaultRegistry
}

// Register registers a configuration entry.
func Register[T any](name string, t *T, doneFunc DoneFunc) {
	defaultRegistry.entries = append(defaultRegistry.entries, registryEntry{
		Query:    name,
		Value:    t,
		DoneFunc: doneFunc,
	})
}

// Decode decodes configuration data into registered entries.
func (r *registry) Decode(ctx context.Context, d data.Data) (err error) {
	defer errors.Handle(&err)

	dec := &data.Decoder{
		TagName:       TagName,
		NameConverter: data.Uncapitalize,
	}

	for _, entry := range r.entries {
		errors.Check(dec.DecodeQuery(d, entry.Query, entry.Value))

		if entry.DoneFunc != nil {
			errors.Check(entry.DoneFunc(ctx))
		}
	}

	return
}
