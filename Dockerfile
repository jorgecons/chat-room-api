# Use an official Golang image to build the app
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download all dependencies
COPY go.mod go.sum ./

# Download all the dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build -o ./cmd/main ./cmd

# Start a new stage from a smaller base image to copy the build artifact
FROM gcr.io/distroless/base

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the builder image
COPY --from=builder /app/cmd/main .
COPY static /root/static

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
