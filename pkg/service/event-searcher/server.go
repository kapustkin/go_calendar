package searcher

import (
	"fmt"
	"log"

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
	log.Printf("connection to kafka established")
	// Get messages from GRPC
	events, err := grpcConn.GetEventsForSend()
	if err != nil {
		return fmt.Errorf("getEventsForSendRequest failed: %v", err.Error())
	}
	log.Printf("received %v events for sending", len(events))
	// Send messages to kafka
	for _, event := range events {
		err = kafkaConn.AddMessage(fmt.Sprintf("%v %v %v", event.Date, event.User, event.Message))
		if err != nil {
			return fmt.Errorf("sending message to kafka failed: %v", err.Error())
		}
		log.Printf("%v - sended to kafka", event.UUID)
		res, err := grpcConn.SetEventAsSended(event.UUID)
		if err != nil {
			return fmt.Errorf("setEventAsSended failed: %v", err.Error())
		}
		if !res {
			return fmt.Errorf("update status for uuid %v failed", event.UUID)
		}
	}

	err = kafkaConn.Close()
	if err != nil {
		return fmt.Errorf("failed close connection to kafka: %v", err.Error())
	}
	log.Printf("success")
	return nil
}
