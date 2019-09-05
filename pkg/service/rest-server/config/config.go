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
}

// InitConfig initial config
func InitConfig() *Config {
	cfg := Config{}
	flag.IntVar(&cfg.Port, "port", 5000, "application port")
	flag.StringVar(&cfg.Host, "host", "localhost", "application host")
	flag.IntVar(&cfg.Logging, "log", 0, "application logger. 0 - Disable, 1 - Standart, 2 - Verbose json")
	flag.Parse()
	log.Printf("Initital app with config %v", cfg)
	return &cfg
}
