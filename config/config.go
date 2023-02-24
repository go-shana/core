// Package config provides a way to load config data from a config file.
package config

import (
	"context"

	"github.com/go-shana/core/errors"
	"github.com/go-shana/core/initer"
	"github.com/go-shana/core/internal/config"
	"github.com/go-shana/core/validator"
)

// New creates a new instance of T which is initialized with the config data matching the query.
//
// The query is a query string that can be used to find the config data in depth.
// The query syntax can be found in `Data.Query` in the package github.com/go-shana/core/data.
//
// Note that the returned value is fully initialized after the config file is loaded.
// Don't use it before that.
//
// Suppose we have a YAML config file like following:
//
//	my:
//	  config:
//	    name: Shana
//	    age: 16
//
// Here is a sample to load a struct from the config:
//
//	type MyConfig struct {
//	    Name string
//	    Age  int
//	    Role string // Role is absent in config. We'll initialize it later.
//	}
//
//	// myConfig is a global variable.
//	var myConfig = config.New[MyConfig]("my.config")
//
// We can provide optional initialization and validation functions for the config.
//
//	// Validate validates the config automatically.
//	func (c *MyConfig) Validate(ctx context.Context) {
//	    errors.Assert(c.Name != "", c.Age > 0)
//	}
//
//	// Init initializes the config automatically.
//	func (c *MyConfig) Init(ctx context.Context) {
//	    c.Role = "Warrior"
//	}
func New[T any](query string) *T {
	var t T
	config.Register(query, &t, func(ctx context.Context) (err error) {
		defer errors.Handle(&err)
		errors.Check(validator.Validate(ctx, &t))
		errors.Check(initer.Init(ctx, &t))
		return
	})
	return &t
}
