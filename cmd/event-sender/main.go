package main

import (
	"log"

	server "github.com/kapustkin/go_calendar/pkg/service/event-sender"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Fatalf("application exception: %v", err.Error())
	}
}
