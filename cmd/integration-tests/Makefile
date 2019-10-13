all: format lint test 

format:
	gofmt -w .

test:
	go test -v -cover ./...

lint:
	$(shell go env GOPATH)/bin/golangci-lint run --enable-all