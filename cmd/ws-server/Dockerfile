# Use official Go image as base
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files from root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code since we need access to the cmd structure
COPY . .

# Build WebSocket server from the specific directory with explicit architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ws-server ./cmd/ws-server

# Use minimal image for runtime
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/ws-server .

# Expose port
EXPOSE 8081

# Run WebSocket server
CMD ["./ws-server"] 