package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	flag "github.com/spf13/pflag"
)

const ENV_PREFIX = "EVENT_SENDER"

// Config app configuration
type Config struct {
	KafkaConnection string `envconfig:"HOST"`
	KafkaTopic      string `envconfig:"TOPIC"`
	Logging         int    `envconfig:"LOGGER"`
}

// InitConfig initial config
func InitConfig() *Config {
	cfg := Config{}
	flag.StringVarP(&cfg.KafkaConnection, "host", "h", "localhost:9092",
		"kafka connection. Default 'localhost:9092'")
	flag.StringVarP(&cfg.KafkaTopic, "topic", "t", "calendar_eventsForSend",
		"kafka topic. Default 'calendar_eventsForSend'")
	flag.IntVarP(&cfg.Logging, "log", "l", 0,
		"application logger. 0 - Disable, 1 - Standart")
	flag.Parse()

	err := envconfig.Process(ENV_PREFIX, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
