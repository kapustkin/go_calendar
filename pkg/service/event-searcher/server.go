package searcher

import (
	"context"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func Execute(topic string, partition int, connString string) {

	conn, err := kafka.DialLeader(context.Background(), "tcp", connString, topic, partition)

	if err != nil {
		log.Fatal(err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}
