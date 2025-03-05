# Use an official Golang image to build the bot
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /bot

# Copy the go.mod and go.sum files to download dependencies
COPY ../../go.mod ../../go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the entire bot module
COPY . .

# Build the bot binary
RUN go build -o /bot/bot ./cmd/main.go

# Use a minimal base image for the runtime
FROM gcr.io/distroless/base

# Set working directory inside the container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /bot/bot .

# Command to run the bot
CMD ["./bot"]
