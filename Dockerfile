# Build stage
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Install build dependencies for SQLite (cached layer)
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./

# Download dependencies (this layer will be cached if go.mod/go.sum don't change)
RUN go mod download

# Install swag (cached if go.mod doesn't change)
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code (only invalidates cache when code changes)
COPY . .

# Generate swagger docs
RUN /go/bin/swag init -g main.go --output ./docs

# Build the application with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o loveapp-backend .

# Final stage - usar alpine para soportar SQLite
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies (cached layer)
RUN apk add --no-cache sqlite-libs ca-certificates

# Copy binary and docs from builder
COPY --from=builder /app/loveapp-backend .
COPY --from=builder /app/docs ./docs

# Create data directory for SQLite database
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Run the application
CMD ["./loveapp-backend"]