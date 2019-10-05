package searcher

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/kapustkin/go_calendar/pkg/logger"
	"github.com/kapustkin/go_calendar/pkg/service/event-searcher/config"
	"github.com/kapustkin/go_calendar/pkg/service/event-searcher/grpc"
	"github.com/kapustkin/go_calendar/pkg/service/event-searcher/kafka"
)

// Run process
func Run() error {
	logger.Init("grpc-server", "0.0.1")
	log.Info("starting app...")
	conf := config.InitConfig()
	log.Infof("use config: %v", conf)
	for {
		err := execute(conf)
		if err != nil {
			return err
		}

		time.Sleep(60 * time.Second)
	}
}

func execute(c *config.Config) error {
	// Init connection to grpc
	grpcConn := grpc.Init(c)

	// Init connection to kafka
	kafkaConn, err := kafka.Init(c)
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("failed connect to kafka: %v", err.Error())
	}
	log.Info("connection to kafka established")
	// Get messages from GRPC
	events, err := grpcConn.GetEventsForSend()
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("getEventsForSendRequest failed: %v", err.Error())
	}
	log.Infof("received %v events for sending", len(events))
	// Send messages to kafka
	for _, event := range events {
		err = kafkaConn.AddMessage(fmt.Sprintf("%v %v %v", event.Date, event.User, event.Message))
		if err != nil {
			log.Error(err.Error())
			return fmt.Errorf("sending message to kafka failed: %v", err.Error())
		}
		log.Infof("%v - sended to kafka", event.UUID)
		res, err := grpcConn.SetEventAsSended(event.UUID)
		if err != nil {
			log.Error(err.Error())
			return fmt.Errorf("setEventAsSended failed: %v", err.Error())
		}
		if !res {
			log.Warnf("update status for uuid %v failed", event.UUID)
			return fmt.Errorf("update status for uuid %v failed", event.UUID)
		}
	}

	err = kafkaConn.Close()
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("failed close connection to kafka: %v", err.Error())
	}
	log.Infof("success")
	return nil
}
