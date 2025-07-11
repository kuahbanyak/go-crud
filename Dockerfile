FROM golang:1.24-alpine

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata wget

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application from cmd directory
RUN go build -o go-crud ./cmd/

# Make sure the binary is executable
RUN chmod +x ./go-crud

# Expose ports
EXPOSE 8080
EXPOSE 8443


# Command to run the application
CMD ["./go-crud"]