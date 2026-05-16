# Stage 1: Build binary
FROM golang:alpine AS builder

# Install git and other dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/main.go

# Stage 2: Final lightweight image
FROM alpine:latest

# Install tzdata and ca-certificates
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .

# Copy migration files (crucial for auto-migration)
COPY --from=builder /app/migrations ./migrations

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./server"]
