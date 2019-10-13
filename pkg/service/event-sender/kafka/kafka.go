package kafka

import (
	"context"

	"github.com/kapustkin/go_calendar/pkg/service/event-sender/config"
	kafkalib "github.com/segmentio/kafka-go"
)

type Kafka struct {
	conn *kafkalib.Reader
}

func Init(c *config.Config) (*Kafka, error) {
	reader := kafkalib.NewReader(kafkalib.ReaderConfig{
		Brokers:  []string{c.KafkaConnection},
		Topic:    c.KafkaTopic,
		GroupID:  "event-searcher-001",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Kafka{reader}, nil
}

func (k *Kafka) GetMessage(ctx context.Context) ([]byte, error) {
	message, err := k.conn.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}
	err = k.conn.CommitMessages(ctx, message)
	if err != nil {
		return nil, err
	}
	return message.Value, nil
}

func (k *Kafka) Close() error {
	err := k.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
