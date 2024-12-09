# Start with the official Go image for ARM64 architecture
FROM golang:1.23.4-bullseye AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Set environment variables for Go
ENV GOOS=linux \
    GOARCH=arm64

# Build the Go application
RUN go build -o app ./cmd

# Start a new minimal image for runtime
FROM cgr.dev/chainguard/wolfi-base

# Set the working directory inside the runtime container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app .

# Expose the port your application uses
EXPOSE 8080

# Grant execution permissions to the app binary
RUN chmod +x ./app

# Set the entrypoint command
# CMD ["./app"]
