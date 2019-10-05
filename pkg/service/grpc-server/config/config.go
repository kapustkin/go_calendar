package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

const ENV_PREFIX = "GRPC_SERVER"

// Config app configuration
type Config struct {
	Host             string `envconfig:"HOST"`
	StorageType      int    `envconfig:"STORAGE"`
	ConnectionString string `envconfig:"CONN"`
}

// InitConfig initial config
func InitConfig() *Config {
	cfg := Config{}
	flag.StringVarP(&cfg.Host, "host", "h", "localhost:5000", "application host")
	flag.IntVarP(&cfg.StorageType, "storage", "s", 1, "application storage. 0 - inmemory, 1 - posgres")
	flag.StringVarP(&cfg.ConnectionString,
		"connection", "c",
		"postgres://log:pass@localhost/ms_calendar?sslmode=disable",
		"connection string for storage")
	flag.Parse()

	err := envconfig.Process(ENV_PREFIX, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
