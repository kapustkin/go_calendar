package config

import (
	"log"

	flag "github.com/spf13/pflag"
)

// Config app configuration
type Config struct {
	Port             int
	Host             string
	StorageType      int
	ConnectionString string
}

// InitConfig initial config
func InitConfig() *Config {
	cfg := Config{}
	flag.IntVarP(&cfg.Port, "port", "p", 5900, "application port")
	flag.StringVarP(&cfg.Host, "host", "h", "localhost", "application host")
	flag.IntVarP(&cfg.StorageType, "storage", "s", 1, "application storage. 0 - inmemory, 1 - posgres")
	flag.StringVarP(&cfg.ConnectionString, "connection", "c", "postgres://postgres:password@localhost/ms_calendar?sslmode=disable", "connection string for storage")
	flag.Parse()
	log.Printf("Initital app with config %v", cfg)
	return &cfg
}
