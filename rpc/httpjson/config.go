package httpjson

import (
	"context"
	"net"

	"github.com/go-shana/core/errors"
)

const defaultPort = 9696

// Config for http server.
type Config struct {
	IP        string `shana:"ip"`   // The IP to bind. If it's not set, all IPs are bound.
	Port      int    `shana:"port"` // The port to listen.
	PkgPrefix string `shana:"-"`    // Filter all exported routes by package prefix.
}

// Validate validates the config.
func (c *Config) Validate(ctx context.Context) {
	if c.IP != "" {
		if ip := net.ParseIP(c.IP); ip == nil {
			errors.Throwf("httpjson: invalid IP in config [ip=%v]", ip)
			return
		}
	}

	if c.Port < 0 || c.Port > 65535 {
		errors.Throwf("httpjson: invalid port in config [port=%v]", c.Port)
		return
	}
}

// Init initializes the config and fills zero values with defaults.
func (c *Config) Init(ctx context.Context) {
	if c.Port == 0 {
		c.Port = defaultPort
	}
}
