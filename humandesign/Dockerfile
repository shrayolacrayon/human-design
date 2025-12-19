# Build stage
FROM golang:1.22-bookworm AS builder

# Install Swiss Ephemeris library
RUN apt-get update && apt-get install -y libswe-dev && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy go mod file
COPY go.mod ./

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -o humandesign ./cmd/server/main.go

# Runtime stage
FROM debian:bookworm-slim

# Install Swiss Ephemeris runtime library
RUN apt-get update && apt-get install -y libswe2.0 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/humandesign .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./humandesign"]
