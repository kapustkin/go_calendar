package main

import (
	server "github.com/kapustkin/go_calendar/pkg/service/event-searcher"
)

func main() {
	server.Execute("my-topic", 0, "192.168.1.242:9092")
}
