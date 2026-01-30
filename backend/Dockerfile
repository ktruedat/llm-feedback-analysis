# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/cmd/main ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/cmd/main .

# Copy config file
COPY config.yaml .

# Copy migrations
COPY migrations ./migrations

# Expose port
EXPOSE 8080

# Run the application
# Default config path is config.yaml, but can be overridden via CMD or docker-compose
ENTRYPOINT ["./main"]
CMD ["-cfg", "config.yaml"]
