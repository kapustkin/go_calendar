package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	flag "github.com/spf13/pflag"
)

const ENV_PREFIX = "REST_SERVER"

// Config app configuration
type Config struct {
	Host    string `envconfig:"HOST"`
	Logging int    `envconfig:"LOGGER"`
	GRPC    string `envconfig:"GRPC_HOST"`
}

// InitConfig initial config
func InitConfig() *Config {
	cfg := Config{}
	flag.StringVarP(&cfg.Host, "host", "h", "localhost:5000", "application host")
	flag.IntVarP(&cfg.Logging, "logger", "l", 1, "application logger. 0 - Disable, 1 - Standart, 2 - Verbose json")
	flag.StringVarP(&cfg.GRPC, "grpc-host", "g", "localhost:5900", "grpc server application address")
	flag.Parse()

	err := envconfig.Process(ENV_PREFIX, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
