
# Go Clean Architecture CRUD API

A RESTful API built with Go implementing Clean Architecture principles, featuring Product and Account management with PostgreSQL database and Swagger documentation.

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org)
[![Gin Framework](https://img.shields.io/badge/Gin-1.9.1-green.svg)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/GORM-1.25.5-red.svg)](https://gorm.io)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)](https://postgresql.org)
[![Swagger](https://img.shields.io/badge/Swagger-Enabled-orange.svg)](https://swagger.io)

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

### Layers Overview

- **Domain Layer**: Contains business entities and repository interfaces
- **Use Case Layer**: Implements business logic and orchestrates data flow
- **Delivery Layer**: HTTP handlers, middleware, and routing
- **Repository Layer**: Data access implementations
- **Package Layer**: Shared utilities and configurations

## 🚀 Features

- ✅ **Clean Architecture** implementation
- ✅ **RESTful API** with proper HTTP methods
- ✅ **PostgreSQL** database integration
- ✅ **GORM** ORM with auto-migrations
- ✅ **Swagger** API documentation
- ✅ **Middleware** support (CORS, Logging, Auth)
- ✅ **Pagination** for list endpoints
- ✅ **Validation** with struct tags
- ✅ **Error handling** with standardized responses
- ✅ **Docker** support for database
- ✅ **Environment** configuration

## 📋 Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose (for database)
- PostgreSQL (if not using Docker)

## 🛠️ Installation
### 1. Clone the repository
### 2. Install dependencies
### 3. Set up environment variables
Create a `.env` file in the root directory:
### 4. Start PostgreSQL with Docker
### 5. Generate Swagger documentation
### 6. Run the application
The server will start on `http://localhost:8080`

## 📚 API Documentation

### Swagger UI
Access the interactive API documentation at: `http://localhost:8080/swagger/index.html`

### Available Endpoints

#### Products

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/products` | Create a new product |
| GET | `/api/v1/products` | Get all products (paginated) |
| GET | `/api/v1/products/{id}` | Get product by ID |
| PUT | `/api/v1/products/{id}` | Update product |
| DELETE | `/api/v1/products/{id}` | Delete product |
| GET | `/api/v1/products/category/{category}` | Get products by category |

#### Accounts

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/accounts` | Create a new account |
| GET | `/api/v1/accounts/{id}` | Get account by ID |
| PUT | `/api/v1/accounts/{id}` | Update account |
| DELETE | `/api/v1/accounts/{id}` | Delete account |

## 💡 Usage Examples

### Create a Product
### Get All Products with Pagination
### Get Product by ID
### Update a Product
## 📊 Database Schema
### Products Table
## 🏃‍♂️ Development
### Running Tests
### Code Generation
### Adding New Entities

To add a new entity (e.g., `Category`):

1. Create entity struct in `internal/domain/entity/category.go`
2. Create repository interface in `internal/domain/repository/category_repository.go`
3. Implement repository in `internal/repository/postgres/category_repository.go`
4. Create use case in `internal/usecase/category_usecase.go`
5. Create HTTP handler in `internal/delivery/http/handler/category_handler.go`
6. Add routes in `cmd/api/main.go`
7. Update migration in `internal/repository/database/migration.go`

## 🐳 Docker

### Start all services
### Services

- **PostgreSQL**: `localhost:5432`
- **pgAdmin**: `localhost:8081` (admin@admin.com / admin)

### Stop services

## 📝 Response Format

### Success Response
### Error Response

### Paginated Response

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [],
  "pagination": {
    "total": 100,
    "limit": 10,
    "offset": 0,
    "page": 1,    
    "pages": 10
  }
}
```

this README provides:

1. **Clear project overview** with badges and architecture diagram
2. **Comprehensive installation instructions** with step-by-step setup
3. **API documentation** with examples and endpoint descriptions
4. **Usage examples** with curl commands
5. **Database schema** information
6. **Development guidelines** for contributing
7. **Docker setup** instructions
8. **Configuration details** with environment variables
9. **Response format** examples
10. **Contributing guidelines** and support information

The README is structured to help both new developers understand
