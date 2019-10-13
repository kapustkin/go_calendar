package sender

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/kapustkin/go_calendar/pkg/logger"
	"github.com/kapustkin/go_calendar/pkg/service/event-sender/config"
	"github.com/kapustkin/go_calendar/pkg/service/event-sender/kafka"
	"github.com/kapustkin/go_calendar/pkg/service/event-sender/sender"
)

// Run process
func Run() error {
	logger.Init("grpc-server", "0.0.1")
	log.Info("starting app...")
	conf := config.InitConfig()
	log.Infof("use config: %v", conf)
	err := execute(conf)
	if err != nil {
		return err
	}
	return nil
}

func execute(c *config.Config) error {
	// Init connection to kafka
	kafkaConn, err := kafka.Init(c)
	if err != nil {
		log.Errorf(err.Error())
		return fmt.Errorf("failed connect to kafka: %v", err.Error())
	}
	log.Info("connection to kafka established")

	for {
		message, err := kafkaConn.GetMessage(context.Background())
		if err != nil {
			log.Errorf(err.Error())
			return fmt.Errorf("failed GetMessage from kafka: %v", err.Error())
		}
		sender.Send(message)
	}
}
