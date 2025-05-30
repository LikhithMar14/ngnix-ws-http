# Use official Go image as base
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the specific server source code
COPY cmd/http-server/ ./cmd/http-server/

# Build HTTP server from the specific directory
RUN go build -o http-server ./cmd/http-server

FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/http-server .

# Expose port
EXPOSE 8080

# Run HTTP server
CMD ["./http-server"] 