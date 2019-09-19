all: format test lint build

install:
	go build -o bin/mc_calendar

format:
	gofmt -w .

test:
	protoc -I=api/proto/calendarpb/ --go_out=plugins=grpc:pkg/api/v1/ api/proto/calendarpb/calendar.proto	
	go test -cover ./...

lint:
	$(shell go env GOPATH)/bin/golangci-lint run --enable-all

build:
	go build -o=bin/grpc-server ./cmd/grpc-server/
	go build -o=bin/rest-server ./cmd/rest-server/
	go build -o=bin/event-searcher ./cmd/event-searcher/
	go build -o=bin/event-sender ./cmd/event-sender/