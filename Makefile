# Go CRUD API Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=go-crud
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_COMPOSE=docker-compose
DOCKER=docker
IMAGE_NAME=go-crud-api

.PHONY: all build clean test deps docker-build docker-up docker-down docker-restart docker-logs help

all: test build

# Build the application
build:
	@echo "Building application..."
	@go build -o ./bin/api ./cmd/api/main.go
	@echo "Build completed!"

# Clean build files
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf ./bin
	@go clean

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Download dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Run the application locally
run:
	@echo "Starting application..."
	@go run ./cmd/api/main.go

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker-compose build --no-cache

docker-up:
	@echo "Starting application with Docker Compose..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker Compose services..."
	@docker-compose down

docker-restart:
	@echo "Restarting Docker Compose services..."
	@docker-compose restart

docker-logs:
	@echo "Showing logs for all services..."
	@docker-compose logs -f

docker-logs-api:
	@echo "Showing logs for API service..."
	@docker-compose logs -f api

docker-logs-db:
	@echo "Showing logs for DB service..."
	@docker-compose logs -f sqlserver

docker-clean:
	@echo "Cleaning Docker resources..."
	@docker-compose down -v --remove-orphans
	@docker system prune -f

# Development commands
dev-setup: deps docker-up
	@echo "Setting up development environment..."
	@cp .env.example .env
	@echo "Please update .env file with your configuration"

dev-reset: docker-clean docker-up

# Database migration (if you add migration support)
migrate:
	@echo "Running database migrations..."
	@docker-compose --profile migration up migrate

# Health check
health:
	@echo "Checking application health..."
	curl -f http://localhost:8080/health || exit 1

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/api/main.go

# Production deployment commands
railway-login:
	@echo "Logging into Railway..."
	railway login

railway-deploy:
	@echo "Deploying to Railway..."
	railway up

railway-logs:
	@echo "Showing Railway deployment logs..."
	railway logs

railway-status:
	@echo "Checking Railway service status..."
	railway status

# Production build and test
prod-build:
	@echo "Building production Docker image..."
	@docker build -f Dockerfile.railway -t go-crud-railway .

prod-test:
	@echo "Testing production build..."
	@docker run --rm -p 8080:8080 --env-file .env.railway go-crud-railway &
	@sleep 10
	@curl -f http://localhost:8080/health || exit 1
	@docker stop $$(docker ps -q --filter ancestor=go-crud-railway)

# Code quality checks
lint:
	@echo "Running Go linter..."
	@golangci-lint run ./...

security-scan:
	@echo "Running security scan..."
	@gosec ./...

# Database operations for production
db-migrate-prod:
	@echo "Running production database migration..."
	@go run ./cmd/migrate/main.go --env=production

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean build files"
	@echo "  test          - Run tests"
	@echo "  deps          - Download dependencies"
	@echo "  tidy          - Tidy dependencies"
	@echo "  run           - Run application locally"
	@echo "  docker-build  - Build Docker images"
	@echo "  docker-up     - Start all services"
	@echo "  docker-down   - Stop all services"
	@echo "  docker-restart- Restart all services"
	@echo "  docker-logs   - Show logs for all services"
	@echo "  docker-clean  - Clean Docker resources"
	@echo "  dev-setup     - Setup development environment"
	@echo "  dev-reset     - Reset development environment"
	@echo "  health        - Check application health"
	@echo "  prod-build    - Build production Docker image"
	@echo "  prod-test     - Test production build"
	@echo "  railway-login - Login to Railway"
	@echo "  railway-deploy- Deploy to Railway"
	@echo "  railway-logs  - Show Railway logs"
	@echo "  railway-status- Check Railway status"
	@echo "  lint          - Run code linter"
	@echo "  security-scan - Run security scan"
	@echo "  help          - Show this help message"
