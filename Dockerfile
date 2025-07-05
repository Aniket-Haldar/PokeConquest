# Build stage
FROM golang:1.24.3-alpine as builder

WORKDIR /app

# Install git (if you use Go modules)
RUN apk --no-cache add git

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy the source
COPY . .

# Build the Go app
RUN go build -o server .

# Final image
FROM alpine:latest

WORKDIR /app

# Install certs for HTTPS
RUN apk --no-cache add ca-certificates

# Copy built binary
COPY --from=builder /app/server .

# Copy frontend files
COPY frontend ./frontend

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]
