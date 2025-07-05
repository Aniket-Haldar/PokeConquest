# Stage 1: Build Go binary
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary
RUN go build -o server .

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app

# Install certificates (needed for HTTPS calls in Go)
RUN apk --no-cache add ca-certificates

# Copy .env file into image
COPY .env /app/.env


# Copy built binary from builder stage
COPY --from=builder /app/server .

# Copy frontend static files
COPY frontend ./frontend

# Expose API port
EXPOSE 8080

# Run the Go server
CMD ["./server"]
