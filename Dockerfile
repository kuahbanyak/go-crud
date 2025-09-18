# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO disabled for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api ./cmd/api/main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/bin/api .

# Copy config files
COPY --from=builder /app/configs ./configs

# Expose port
EXPOSE 8080

# Run the application
CMD ["./api"]
