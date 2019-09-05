package main

import (
	"log"
	"os"

	rest "github.com/kapustkin/go_calendar/pkg/service/rest-server"
)

func main() {
	err := rest.Run(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
