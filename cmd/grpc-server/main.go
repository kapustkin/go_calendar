package main

import (
	"log"

	"github.com/kapustkin/go_calendar/pkg/service/grpc-server"
)

func main() {
	err := grpc.Run()
	if err != nil {
		log.Fatal(err)
	}
}
