package searcher

import (
	"fmt"

	"github.com/kapustkin/go_calendar/pkg/service/event-searcher/config"
	"github.com/kapustkin/go_calendar/pkg/service/event-searcher/grpc"
	"github.com/kapustkin/go_calendar/pkg/service/event-searcher/kafka"
)

// Run process
func Run() error {
	c := config.InitConfig()
	err := execute(c)
	if err != nil {
		return err
	}
	return nil
}

func execute(c *config.Config) error {
	// Init connection to grpc
	grpcConn := grpc.Init(c)

	// Init connection to kafka
	kafkaConn, err := kafka.Init(c)
	if err != nil {
		return fmt.Errorf("failed connect to kafka: %v", err.Error())
	}
	// Get messages from GRPC
	events, err := grpcConn.GetEventsForNotify()
	if err != nil {
		return fmt.Errorf("grpc GetEventsForSendRequest error: %v", err.Error())
	}

	// Send messages to kafka
	for _, event := range events {
		err = kafkaConn.AddMessage(fmt.Sprintf("%v %v %v", event.Date, event.User, event.Message))
		if err != nil {
			return fmt.Errorf("error sending message to kafka: %v", err.Error())
		}
	}

	err = kafkaConn.Close()
	if err != nil {
		return fmt.Errorf("failed close connection to kafka: %v", err.Error())
	}

	return nil
}
