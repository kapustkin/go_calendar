all: format lint test 

format:
	gofmt -w .

test:
	go test -cover ./...

lint:
	$(shell go env GOPATH)/bin/golangci-lint run --enable-all