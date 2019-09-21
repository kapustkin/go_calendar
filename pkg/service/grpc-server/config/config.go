package config

import (
	"log"

	flag "github.com/spf13/pflag"
	viper "github.com/spf13/viper"
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
	flag.StringVarP(&cfg.ConnectionString,
		"connection", "c",
		"postgres://log:pass@localhost/ms_calendar?sslmode=disable",
		"connection string for storage")
	flag.Parse()

	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Fatal(err)
	}

	viper.SetEnvPrefix("cal_grpc")

	_ = viper.BindEnv("port")
	_ = viper.BindEnv("host")
	_ = viper.BindEnv("storage")
	_ = viper.BindEnv("conn")

	cfg.Port = viper.GetInt("port")
	cfg.Host = viper.GetString("host")
	cfg.StorageType = viper.GetInt("storage")
	cfg.ConnectionString = viper.GetString("conn")
	return &cfg
}
