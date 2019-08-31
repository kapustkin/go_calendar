package main

import (
	"log"
	"os"

	"github.com/kapustkin/go_calendar/internal"
)

func main() {
	err := internal.Run(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
