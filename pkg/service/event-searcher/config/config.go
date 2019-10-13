package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	flag "github.com/spf13/pflag"
)

const envPrefix = "EVENT_SEARCHER"

// Config app configuration
type Config struct {
	KafkaConnection  string `envconfig:"HOST"`
	KafkaTopic       string `envconfig:"TOPIC"`
	KafkaPartition   int    `envconfig:"PARTITION"`
	GrpcConnection   string `envconfig:"GRPC_HOST"`
	EventSearchDelay int64  `envconfig:"DELAY"`
	Logging          int    `envconfig:"LOGGER"`
}

// InitConfig initial config
func InitConfig() *Config {
	cfg := Config{}
	flag.StringVarP(&cfg.KafkaConnection, "host", "h", "localhost:9092",
		"kafka connection. Default 'localhost:9092'")
	flag.StringVarP(&cfg.KafkaTopic, "topic", "t", "eventsForSend",
		"kafka topic. Default 'eventsForSend'")
	flag.IntVarP(&cfg.KafkaPartition, "partiotion", "p", 0,
		"kafka partiotion. Default '0'")
	flag.StringVarP(&cfg.GrpcConnection, "grpc", "g", "localhost:5900",
		"grpc server application address")
	flag.Int64VarP(&cfg.EventSearchDelay, "delay", "d", 10,
		"application event search delay. Default '10'")
	flag.IntVarP(&cfg.Logging, "log", "l", 0,
		"application logger. 0 - Disable, 1 - Standart")
	flag.Parse()

	err := envconfig.Process(envPrefix, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
