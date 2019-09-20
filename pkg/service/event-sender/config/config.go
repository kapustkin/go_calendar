package config

import (
	"log"

	flag "github.com/spf13/pflag"
)

// Config app configuration
type Config struct {
	KafkaConnection string
	KafkaTopic      string
	KafkaPartition  int
	Logging         int
}

// InitConfig initial config
func InitConfig() *Config {
	cfg := Config{}
	flag.StringVarP(&cfg.KafkaConnection, "host", "h", "192.168.1.242:9092",
		"kafka connection. Default '192.168.1.242:9092'")
	flag.StringVarP(&cfg.KafkaTopic, "topic", "t", "calendar_eventsForSend",
		"kafka topic. Default 'calendar_eventsForSend'")
	flag.IntVarP(&cfg.KafkaPartition, "partiotion", "p", 0,
		"kafka partiotion. Default '0'")
	flag.IntVarP(&cfg.Logging, "log", "l", 0,
		"application logger. 0 - Disable, 1 - Standart")

	flag.Parse()
	log.Printf("Initital app with config %v", cfg)
	return &cfg
}
