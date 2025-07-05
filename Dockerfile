# Build stage
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Install git & certs
RUN apk --no-cache add git ca-certificates

# Copy go.mod and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN go build -o server .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install certs for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary & frontend files
COPY --from=builder /app/server .
COPY frontend ./frontend

# Expose port for Render
EXPOSE 8080

# Start server
CMD ["./server"]
