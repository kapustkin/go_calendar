package config

import (
	"log"

	flag "github.com/spf13/pflag"
)

// Config app configuration
type Config struct {
	Port    int
	Host    string
	Logging int
	GRPC    string
}

// InitConfig initial config
func InitConfig() *Config {
	cfg := Config{}
	flag.IntVarP(&cfg.Port, "port", "p", 5000, "application port")
	flag.StringVarP(&cfg.Host, "host", "h", "localhost", "application host")
	flag.IntVarP(&cfg.Logging, "log", "l", 0, "application logger. 0 - Disable, 1 - Standart, 2 - Verbose json")
	flag.StringVarP(&cfg.GRPC, "grpc", "g", "localhost:5900", "grpc server application address")
	flag.Parse()
	log.Printf("Initital app with config %v", cfg)
	return &cfg
}
