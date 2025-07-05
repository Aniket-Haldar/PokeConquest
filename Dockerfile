# Build stage
FROM golang:1.24.3-alpine as builder

WORKDIR /app

# Install git & ca-certs
RUN apk --no-cache add git ca-certificates

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the app
RUN go build -o server .

# Final image
FROM alpine:latest

WORKDIR /app

# Copy certs & frontend
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/server .
COPY frontend ./frontend

# Expose port (Render uses $PORT)
EXPOSE 8080

# Run the server
CMD ["./server"]
