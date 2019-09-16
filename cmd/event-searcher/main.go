package main

import (
	"context"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "192.168.1.242:9092", topic, partition)

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
