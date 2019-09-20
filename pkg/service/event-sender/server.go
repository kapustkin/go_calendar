package sender

import (
	"context"
	"fmt"
	"log"

	"github.com/kapustkin/go_calendar/pkg/service/event-sender/config"
	"github.com/kapustkin/go_calendar/pkg/service/event-sender/kafka"
	"github.com/kapustkin/go_calendar/pkg/service/event-sender/sender"
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
	// Init connection to kafka
	kafkaConn, err := kafka.Init(c)
	if err != nil {
		return fmt.Errorf("failed connect to kafka: %v", err.Error())
	}
	log.Printf("connection to kafka established")
	ctx := context.Background()
	for {
		message, err := kafkaConn.GetMessage(ctx)
		if err != nil {
			return fmt.Errorf("failed GetMessage from kafka: %v", err.Error())
		}
		sender.Send(message)
		//log.Printf("recieved user event: %v", string(message))
	}
}
