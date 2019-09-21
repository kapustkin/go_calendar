# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Alexandr Kapustkin <a.kapustkin@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Build the Go app
RUN go build -o main ./cmd/grpc-server/

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]