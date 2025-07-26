FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# Pre-copy/cache go.mod for pre-downloading dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/adserver ./cmd/server

# Create a minimal production image
FROM alpine:3.18

WORKDIR /app

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/bin/adserver /app/adserver

# Create a non-root user to run the application
RUN adduser -D appuser && \
    chown -R appuser:appuser /app

USER appuser

# Expose the application port
EXPOSE 8080

# Set the entry point
ENTRYPOINT ["/app/adserver"]