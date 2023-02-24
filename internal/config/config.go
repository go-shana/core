package config

import (
	"context"
	"os"

	"github.com/go-shana/core/data"
	"github.com/go-shana/core/errors"
	"gopkg.in/yaml.v3"
)

var errInvalidConfigFile = errors.New("config: invalid config file")

// Config is a parsed config file.
type Config struct {
	data data.Data
}

// New creates a new config.
func New() *Config {
	return &Config{}
}

// Load parses a config file and merges parsed data with existing data.
func (c *Config) Load(ctx context.Context, filename string) (err error) {
	defer errors.Handle(&err)

	file := errors.Check1(os.Open(filename))
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	raw := data.RawData{}
	errInvalidConfigFile.Check(decoder.Decode(&raw))
	d := data.Make(raw)
	data.MergeTo(&c.data, d)

	return
}

// Data returns parsed data.
func (c *Config) Data() data.Data {
	return c.data
}
