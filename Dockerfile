# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Run go mod tidy to ensure dependencies are up to date
RUN go mod tidy

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o geolocation-service ./cmd/api

# Final stage - using alpine
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/geolocation-service .

# Copy any additional files needed at runtime
COPY --from=builder /app/scripts/migrations/ /app/migrations/

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./geolocation-service"]