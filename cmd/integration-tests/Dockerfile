FROM golang:1.12.6

# Set the Current Working Directory inside the container
WORKDIR /go/bin
# Copy files
COPY . .
# Update packages
RUN curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.20.0 && go get -t -d -v
# Build the Go app
CMD ["make", "all"]
