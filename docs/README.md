# Go CRUD API - Clean Architecture

A Go-based REST API built with Clean Architecture principles for managing users, vehicles, bookings, and inventory in an automotive service system.

## Project Structure

```
project-name/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── domain/                     # Business entities and rules
│   │   ├── entities/               # Core business entities
│   │   │   ├── user.go
│   │   │   ├── vehicle.go
│   │   │   ├── booking.go
│   │   │   ├── invoice.go
│   │   │   └── part.go
│   │   ├── repositories/           # Repository interfaces
│   │   │   ├── user_repository.go
│   │   │   ├── vehicle_repository.go
│   │   │   ├── booking_repository.go
│   │   │   └── inventory_repository.go
│   │   └── services/               # Domain service interfaces
│   │       ├── user_service.go
│   │       └── booking_service.go
│   ├── usecases/                   # Application business rules
│   │   ├── user_usecase.go
│   │   └── booking_usecase.go
│   ├── adapters/
│   │   ├── handlers/               # HTTP handlers (controllers)
│   │   │   ├── http/
│   │   │   │   ├── user_handler.go
│   │   │   │   ├── booking_handler.go
│   │   │   │   └── middleware/
│   │   │   │       ├── auth.go
│   │   │   │       ├── cors.go
│   │   │   │       └── logging.go
│   │   └── repositories/           # Repository implementations
│   │       └── mssql/
│   │           ├── user_repository.go
│   │           └── booking_repository.go
│   ├── infrastructure/             # Framework and drivers
│   │   ├── database/
│   │   │   └── mssql.go
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── logger/
│   │   │   └── logger.go
│   │   └── server/
│   │       └── http.go
│   └── shared/                     # Shared utilities
│       ├── utils/
│       │   └── jwt.go
│       ├── dto/
│       │   └── user_dto.go
│       └── constants/
│           └── constants.go
├── pkg/                           # Public libraries
│   ├── errors/
│   │   └── errors.go
│   └── response/
│       └── response.go
├── configs/                       # Configuration files
│   ├── config.yaml
│   └── config.dev.yaml
├── scripts/                       # Build and deployment scripts
│   ├── build.sh
│   └── migrate.sh
├── .env.example
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

## Key Features

- **Clean Architecture**: Follows Uncle Bob's Clean Architecture principles
- **UUID Primary Keys**: All entities use UUID with SQL Server `uniqueidentifier` type
- **JWT Authentication**: Secure authentication with role-based access control
- **Repository Pattern**: Interface-based repository pattern for data access
- **Dependency Injection**: Proper dependency injection throughout the application
- **Middleware Support**: CORS, logging, authentication middleware
- **Docker Support**: Full Docker and Docker Compose support

## UUID Implementation

All models now use UUID with the following pattern:

```go
type Entity struct {
    ID        uuid.UUID      `gorm:"type:uniqueidentifier;primary_key;default:newid()" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    // ... other fields
}

func (e *Entity) BeforeCreate(tx *gorm.DB) error {
    if e.ID == uuid.Nil {
        e.ID = uuid.New()
    }
    return nil
}
```

## Getting Started

### Prerequisites

- Go 1.21+
- SQL Server (or Docker)
- Redis (optional, for caching)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go-crud
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Update the `.env` file with your database and other configurations.

4. Install dependencies:
```bash
make deps
```

5. Run database migrations:
```bash
make migrate
```

6. Build and run the application:
```bash
make build
make run
```

### Using Docker

1. Start with Docker Compose:
```bash
make docker-run
```

This will start:
- SQL Server database
- Redis cache
- The Go API application

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user

### Users (Protected)
- `GET /api/v1/users/profile` - Get user profile

### Bookings (Protected)
- `POST /api/v1/bookings` - Create a new booking
- `GET /api/v1/bookings/:id` - Get booking by ID
- `PUT /api/v1/bookings/:id/assign-mechanic` - Assign mechanic to booking

## Configuration

The application supports configuration via:
- Environment variables (`.env` file)
- YAML configuration files (`configs/config.yaml`)
- Command line arguments

### Environment Variables

```env
SERVER_PORT=8080
SERVER_HOST=localhost
DB_HOST=localhost
DB_PORT=1433
DB_USER=sa
DB_PASSWORD=your_password
DB_DATABASE=go_crud
JWT_SECRET=your-jwt-secret
JWT_EXPIRATION=24
```

## Development

### Make Commands

- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make migrate` - Run database migrations
- `make deps` - Install dependencies
- `make docker-build` - Build Docker image
- `make docker-run` - Run with Docker Compose
- `make lint` - Run linter
- `make fmt` - Format code

### Project Layers

1. **Domain Layer** (`internal/domain/`): Contains business entities, repository interfaces, and domain services
2. **Use Cases Layer** (`internal/usecases/`): Contains application-specific business logic
3. **Adapters Layer** (`internal/adapters/`): Contains handlers, repository implementations, and external service adapters
4. **Infrastructure Layer** (`internal/infrastructure/`): Contains framework-specific code (database, server, config)
5. **Shared Layer** (`internal/shared/`): Contains utilities, DTOs, and constants shared across layers

### Clean Architecture Benefits

- **Independence**: Each layer is independent and can be tested in isolation
- **Testability**: Business logic can be tested without UI, database, or external services
- **Flexibility**: Easy to swap implementations (database, UI, external services)
- **Maintainability**: Clear separation of concerns makes the code easier to maintain

## Testing

Run all tests:
```bash
make test
```

Run specific test packages:
```bash
go test ./internal/usecases/...
go test ./internal/adapters/repositories/...
```

## Deployment

### Production Build

```bash
make build
```

### Docker Deployment

```bash
make docker-build
docker tag go-crud-api:latest your-registry/go-crud-api:latest
docker push your-registry/go-crud-api:latest
```

## Contributing

1. Follow clean architecture principles
2. Write tests for new features
3. Use conventional commit messages
4. Ensure all lints pass with `make lint`
5. Format code with `make fmt`

## License

This project is licensed under the MIT License.
