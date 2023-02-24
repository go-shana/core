package httpjson

import (
	"context"
)

// Config for http server.
type Config struct {
	Debug     bool   // Debug mode.
	Addr      string // The address to start to listen.
	PkgPrefix string // Filter all exported routes by package prefix.
}

// Init initializes the config and fills zero values with defaults.
func (c *Config) Init(ctx context.Context) (err error) {
	if c.Addr == "" {
		c.Addr = ":9696"
	}

	return
}
