package launcher

import (
	"flag"
	"fmt"
	"os"
)

const defaultConfig = "shana.yaml"
const extConfigEnv = "SHANA_CONFIG_EXT"

var (
	flagVersion = flag.Bool("version", false, "Show service version and exit")
	flagConfig  = flag.String("config", defaultConfig, "Load the config `filename`. Default filename is 'shana.yaml'.")
)

type cliConfig struct {
	MainConfig string
	ExtConfig  string
}

func parseFlags() *cliConfig {
	flag.Parse()

	if *flagVersion {
		fmt.Fprintln(os.Stderr, metaData["version"])
		os.Exit(1)
		panic("never reach here")
	}

	main := *flagConfig

	if main == "" {
		main = defaultConfig
	}

	ext, _ := os.LookupEnv(extConfigEnv)

	return &cliConfig{
		MainConfig: main,
		ExtConfig:  ext,
	}
}
