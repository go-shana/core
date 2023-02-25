package global

import "github.com/go-shana/core/config"

// Config is the configuration for global microservice.
type Config struct {
	Debug bool
}

var (
	defaultConfig = config.New[Config]("shana")
)

// Debug returns global debug flag.
func Debug() bool {
	return defaultConfig.Debug
}
