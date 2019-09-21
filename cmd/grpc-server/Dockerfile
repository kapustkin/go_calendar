# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Alexandr Kapustkin <a.kapustkin@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . .
# Build the Go app
RUN go build -o rest-server ./...

# Expose port 5900 to the outside world
EXPOSE 5900

# Command to run the executable
CMD ["./rest-server"]