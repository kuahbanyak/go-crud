# Go CRUD API Makefile

.PHONY: build run test clean migrate docker-build docker-run

# Build the application
build:
	@echo "Building application..."
	@go build -o ./bin/api ./cmd/api/main.go
	@echo "Build completed!"

# Run the application
run:
	@echo "Starting application..."
	@go run ./cmd/api/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf ./bin
	@go clean

# Run database migrations
migrate:
	@echo "Running database migrations..."
	@./scripts/migrate.sh

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t go-crud-api .

# Run with Docker Compose
docker-run:
	@echo "Starting application with Docker Compose..."
	@docker-compose up -d

# Stop Docker Compose
docker-stop:
	@echo "Stopping Docker Compose services..."
	@docker-compose down

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	@cp .env.example .env
	@echo "Please update .env file with your configuration"

# Lint code
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
