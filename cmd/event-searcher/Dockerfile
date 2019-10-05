# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Alexandr Kapustkin <a.kapustkin@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /go/bin
# Copy files
COPY . .
# Update packages
RUN go get -d -v
# Build the Go app
RUN CGO_ENABLED=0 go build -o event-searcher ./...

# STEP 2 build a small image
# start from scratch
FROM alpine:latest
# Copy our static executable
COPY --from=builder /go/bin/event-searcher /go/bin/event-searcher
ENTRYPOINT ["/go/bin/event-searcher"]