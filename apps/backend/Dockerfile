# Build Stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY apps/backend/go.mod apps/backend/go.sum ./apps/backend/
WORKDIR /app/apps/backend

# Download dependencies
RUN go mod download

# Copy the rest of the backend source
COPY apps/backend/ .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /api ./cmd/api/main.go

# Final Stage
FROM alpine:3.19

# Add CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /api /app/api

# Expose the API port
EXPOSE 8080

# Command to run the application
CMD ["/app/api"]
