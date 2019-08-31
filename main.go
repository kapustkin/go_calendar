package main

import (
	"log"
	"os"

	"github.com/kapustkin/go_ms_template/internal"
)

func main() {
	err := internal.Run(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
