package main

import (
	"log"

	"github.com/kapustkin/go_calendar/pkg/service/grpc-server"
)

func main() {
	err := grpc.Run("0.0.0.0:5900")
	if err != nil {
		log.Fatal(err)
	}
}
