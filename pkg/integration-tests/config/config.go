package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	flag "github.com/spf13/pflag"
)

const envPrefix = "INTEGRATION_TESTS"

// Config app configuration
type Config struct {
	KafkaConnection string `envconfig:"KAFKA_HOST"`
	KafkaTopic      string `envconfig:"KAFKA_TOPIC"`
	RestServer      string `envconfig:"REST_SERVER"`
}

// InitConfig initial config
func InitConfig() *Config {
	cfg := Config{}
	flag.StringVarP(&cfg.KafkaConnection, "host", "h", "localhost:9092",
		"kafka connection. Default 'localhost:9092'")
	flag.StringVarP(&cfg.KafkaTopic, "topic", "t", "calendar_eventsForSend",
		"kafka topic. Default 'calendar_eventsForSend'")
	flag.StringVarP(&cfg.RestServer, "rest", "r", "localhost:5000",
		"rest server application address")
	flag.Parse()

	err := envconfig.Process(envPrefix, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
