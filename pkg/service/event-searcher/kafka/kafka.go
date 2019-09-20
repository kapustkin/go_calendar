package kafka

import (
	"context"
	"time"

	"github.com/kapustkin/go_calendar/pkg/service/event-searcher/config"
	kafkalib "github.com/segmentio/kafka-go"
)

type Kafka struct {
	conn *kafkalib.Conn
}

func Init(c *config.Config) (*Kafka, error) {
	conn, err := kafkalib.DialLeader(context.Background(), "tcp", c.KafkaConnection, c.KafkaTopic, c.KafkaPartition)
	if err != nil {
		return nil, err
	}
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return nil, err
	}
	return &Kafka{conn}, nil
}

func (k *Kafka) AddMessage(message string) error {
	_, err := k.conn.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func (k *Kafka) Close() error {
	err := k.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
