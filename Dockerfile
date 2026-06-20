# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the server application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mcp-tududi ./cmd/server

# Build the CLI application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mcp-cli ./cmd/cli

# Production stage - Server
FROM alpine:latest AS server

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/mcp-tududi .

# Set environment variables
ENV HTTP_PORT=8080 \
    LOG_LEVEL=info

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the server
CMD ["./mcp-tududi"]

# CLI stage
FROM alpine:latest AS cli

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/mcp-cli .

# Default entrypoint
ENTRYPOINT ["./mcp-cli"]
