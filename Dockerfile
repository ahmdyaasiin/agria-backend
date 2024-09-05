# Stage 1: Build the Go binary
FROM golang:1.22.6 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download go module dependencies
RUN go mod download

# Copy the rest of the files
COPY . .

COPY .env .

# Build the Go application with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/app

# Stage 2: Create a lightweight image to run the binary
FROM alpine:latest

# Install necessary dependencies for running the binary
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Ensure the binary has execution permissions
RUN chmod +x main

# Set the command to run the binary
CMD ["./main"]
