# Car Maintenance Service Management System

A comprehensive RESTful API service for managing car maintenance operations, built with Go using Clean Architecture principles. This system handles customer queuing, vehicle management, maintenance tracking, and service progress monitoring.

## 🚀 Features

### Core Functionality
- **User Management**: Registration, authentication, and role-based access control (Customer, Mechanic, Admin)
- **Vehicle Management**: Track customer vehicles with detailed information
- **Queue Management**: Digital waiting list system with real-time status updates
- **Maintenance Items**: Track maintenance tasks with approval workflow
- **Service Progress**: Real-time service progress tracking for customers
- **Product Management**: Inventory management for parts and products
- **Settings Management**: Configurable shop settings and business hours
- **Notification System**: Event-driven notifications via RabbitMQ

### Technical Features
- Clean Architecture with clear separation of concerns
- RESTful API with consistent response format
- JWT-based authentication
- Role-based authorization (Admin, Mechanic, Customer)
- Request ID tracking for debugging
- Rate limiting and request size validation
- CORS support
- Comprehensive logging
- Database migrations with GORM
- Background job scheduling
- Event-driven architecture with RabbitMQ
- Docker containerization

## 📋 Prerequisites

- **Go** 1.24 or higher
- **SQL Server** 2022 or compatible
- **RabbitMQ** (for message queuing)
- **Docker & Docker Compose** (optional, for containerized deployment)

## 🛠️ Tech Stack

### Backend
- **Go** 1.24
- **Gorilla Mux** - HTTP router
- **GORM** - ORM library
- **JWT** - Authentication
- **Bcrypt** - Password hashing
- **RabbitMQ** - Message broker
- **GoCron** - Job scheduling

### Database
- **Microsoft SQL Server** 2022

### Infrastructure
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration

## 📁 Project Structure

```
go-crud/
├── cmd/
│   ├── api/                    # API server entry point
│   └── worker/                 # Background worker entry point
├── configs/                    # Configuration files
│   ├── config.yaml
│   └── config.prod.yaml
├── internal/
│   ├── adapters/              # External adapters
│   │   ├── handlers/          # HTTP handlers
│   │   │   └── http/
│   │   │       ├── middleware/  # HTTP middleware
│   │   │       ├── *_handler.go # Route handlers
│   │   ├── repositories/      # Data repositories
│   │   │   └── mssql/
│   │   └── external/          # External service clients
│   │       ├── email/
│   │       └── payment/
│   ├── domain/                # Domain layer
│   │   ├── entities/          # Domain entities
│   │   ├── repositories/      # Repository interfaces
│   │   └── services/          # Domain services
│   ├── infrastructure/        # Infrastructure layer
│   │   ├── config/            # Config management
│   │   ├── database/          # Database connection
│   │   ├── logger/            # Logging
│   │   ├── messaging/         # RabbitMQ integration
│   │   ├── scheduler/         # Job scheduler
│   │   └── server/            # HTTP server
│   ├── shared/                # Shared utilities
│   │   ├── constants/         # Application constants
│   │   ├── dto/               # Data transfer objects
│   │   ├── types/             # Custom types
│   │   └── utils/             # Utility functions
│   └── usecases/              # Business logic
├── pkg/                       # Public packages
│   └── response/              # HTTP response utilities
├── tests/                     # Tests
│   ├── mocks/                 # Mock implementations
│   └── unit/                  # Unit tests
├── docker-compose.yml         # Docker composition
├── Dockerfile                 # API container
├── Dockerfile.worker          # Worker container
└── go.mod                     # Go modules
```

## 🚦 Getting Started

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/kuahbanyak/go-crud.git
cd go-crud
```

2. **Copy environment configuration**
```bash
cp .env.example .env
```

3. **Configure environment variables**

Edit `.env` file with your settings:
```env
# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=1433
DB_USER=sa
DB_PASSWORD=YourStrong@Passw0rd
DB_DATABASE=gocrud

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24

# RabbitMQ Configuration
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=admin
RABBITMQ_PASS=rabbitmq_secure_password_123
```

### Running with Docker Compose (Recommended)

```bash
docker-compose up -d
```

This will start:
- API Server on port 8081
- SQL Server on port 1433
- RabbitMQ on port 5672 (Management UI: http://localhost:15672)
- Notification Worker

### Running Locally

1. **Install dependencies**
```bash
go mod download
```

2. **Ensure SQL Server and RabbitMQ are running**

3. **Run database migrations** (GORM auto-migrate will run on startup)

4. **Start the API server**
```bash
go run cmd/api/main.go
```

5. **Start the worker (optional, in separate terminal)**
```bash
go run cmd/worker/main.go
```

The API will be available at `http://localhost:8080`

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Health Check
```http
GET /health
```

### Authentication Endpoints

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe",
  "phone": "08123456789",
  "role": "customer"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "customer"
    }
  }
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
Authorization: Bearer {token}
```

### User Profile Endpoints

#### Get Profile
```http
GET /api/v1/users/profile
Authorization: Bearer {token}
```

#### Update Profile
```http
PUT /api/v1/users/profile
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Jane Doe",
  "phone": "08123456789"
}
```

### Vehicle Management

#### Create Vehicle
```http
POST /api/v1/vehicles
Authorization: Bearer {token}
Content-Type: application/json

{
  "license_plate": "B1234XYZ",
  "brand": "Toyota",
  "model": "Avanza",
  "year": 2022,
  "color": "Silver"
}
```

#### Get My Vehicles
```http
GET /api/v1/vehicles
Authorization: Bearer {token}
```

#### Get Vehicle by ID
```http
GET /api/v1/vehicles/{id}
Authorization: Bearer {token}
```

#### Update Vehicle
```http
PUT /api/v1/vehicles/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "color": "Black",
  "year": 2022
}
```

#### Delete Vehicle
```http
DELETE /api/v1/vehicles/{id}
Authorization: Bearer {token}
```

### Waiting List / Queue Management

#### Take Queue Number
```http
POST /api/v1/waiting-list/take
Authorization: Bearer {token}
Content-Type: application/json

{
  "vehicle_id": "uuid",
  "service_date": "2025-11-15",
  "estimated_time": 120,
  "notes": "Oil change and tire rotation"
}
```

#### Get My Queue
```http
GET /api/v1/waiting-list/my-queue
Authorization: Bearer {token}
```

#### Get Today's Queue
```http
GET /api/v1/waiting-list/today
Authorization: Bearer {token}
```

#### Get Queue by Date
```http
GET /api/v1/waiting-list/date?date=2025-11-15
Authorization: Bearer {token}
```

#### Get Queue by Number
```http
GET /api/v1/waiting-list/number/{number}
Authorization: Bearer {token}
```

#### Check Availability
```http
GET /api/v1/waiting-list/availability?date=2025-11-15
Authorization: Bearer {token}
```

#### Cancel Queue
```http
PUT /api/v1/waiting-list/{id}/cancel
Authorization: Bearer {token}
```

#### Get Service Progress
```http
GET /api/v1/waiting-list/{id}/progress
Authorization: Bearer {token}
```

### Maintenance Items

#### Create Initial Maintenance Items
```http
POST /api/v1/maintenance/waiting-list/{waiting_list_id}/items
Authorization: Bearer {token}
Content-Type: application/json

{
  "items": [
    {
      "category": "Engine",
      "name": "Oil Change",
      "description": "Change engine oil",
      "estimated_cost": 150000,
      "estimated_time": 30
    }
  ]
}
```

#### Get Items by Waiting List
```http
GET /api/v1/maintenance/waiting-list/{waiting_list_id}/items
Authorization: Bearer {token}
```

#### Get Inspection Summary
```http
GET /api/v1/maintenance/waiting-list/{waiting_list_id}/inspection-summary
Authorization: Bearer {token}
```

#### Approve/Reject Maintenance Items
```http
POST /api/v1/maintenance/items/approve
Authorization: Bearer {token}
Content-Type: application/json

{
  "item_ids": ["uuid1", "uuid2"],
  "approve": true,
  "notes": "Approved all items"
}
```

### Products

#### Get All Products
```http
GET /api/v1/products
```

#### Get Product by ID
```http
GET /api/v1/products/{id}
```

### Settings

#### Get Public Settings
```http
GET /api/v1/settings/public
Authorization: Bearer {token}
```

### Admin Endpoints

All admin endpoints require `Authorization: Bearer {admin_token}` and admin role.

#### Waiting List Management
```http
PUT /api/v1/admin/waiting-list/{id}/call       # Call customer
PUT /api/v1/admin/waiting-list/{id}/start      # Start service
PUT /api/v1/admin/waiting-list/{id}/complete   # Complete service
PUT /api/v1/admin/waiting-list/{id}/no-show    # Mark no-show
```

#### Maintenance Items (Mechanic/Admin)
```http
POST /api/v1/admin/maintenance/items/discovered  # Add discovered issue
PUT /api/v1/admin/maintenance/items/{id}         # Update item
PUT /api/v1/admin/maintenance/items/{id}/complete # Complete item
DELETE /api/v1/admin/maintenance/items/{id}      # Delete item
```

#### Product Management
```http
POST /api/v1/admin/products           # Create product
PUT /api/v1/admin/products/{id}       # Update product
PATCH /api/v1/admin/products/{id}/stock # Update stock
DELETE /api/v1/admin/products/{id}    # Delete product
```

#### User Management
```http
GET /api/v1/users                     # Get all users
GET /api/v1/users/{id}                # Get user by ID
PUT /api/v1/users/{id}                # Update user
DELETE /api/v1/users/{id}             # Delete user
```

#### Settings Management
```http
GET /api/v1/admin/settings                    # Get all settings
POST /api/v1/admin/settings                   # Create setting
GET /api/v1/admin/settings/category/{category} # Get by category
GET /api/v1/admin/settings/key/{key}          # Get by key
PUT /api/v1/admin/settings/key/{key}          # Update setting
DELETE /api/v1/admin/settings/{id}            # Delete setting
```

#### Vehicle Management (Admin)
```http
GET /api/v1/admin/vehicles            # Get all vehicles
```

## 🔐 Authentication & Authorization

### Roles
- **Customer**: Can manage their own vehicles, take queue numbers, view progress
- **Mechanic**: Can update maintenance items, mark service progress
- **Admin**: Full access to all operations

### JWT Token
Include the JWT token in the Authorization header:
```
Authorization: Bearer {your-jwt-token}
```

## 📊 Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    ...
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

## 🧪 Testing

### Run Unit Tests
```bash
go test ./tests/unit/... -v
```

### Run All Tests
```bash
go test ./... -v
```

### Run with Coverage
```bash
go test ./... -cover
```

## 🐳 Docker Commands

### Build and Run
```bash
docker-compose up --build
```

### Stop Services
```bash
docker-compose down
```

### View Logs
```bash
docker-compose logs -f go-crud-api
docker-compose logs -f notification-worker
```

### Rebuild Specific Service
```bash
docker-compose up --build go-crud-api
```

## 📦 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | HTTP server port | 8080 |
| GIN_MODE | Gin mode (debug/release) | debug |
| LOG_LEVEL | Logging level | info |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 1433 |
| DB_USER | Database user | sa |
| DB_PASSWORD | Database password | - |
| DB_DATABASE | Database name | gocrud |
| JWT_SECRET | JWT secret key | - |
| JWT_EXPIRATION | JWT expiration (hours) | 24 |
| RABBITMQ_HOST | RabbitMQ host | localhost |
| RABBITMQ_PORT | RabbitMQ port | 5672 |
| RABBITMQ_USER | RabbitMQ username | admin |
| RABBITMQ_PASS | RabbitMQ password | - |

## 🔧 Configuration

### Database Connection
Configure in `.env` or `configs/config.yaml`:
```yaml
database:
  host: "localhost"
  port: "1433"
  user: "sa"
  password: "YourPassword"
  database: "gocrud"
```

### JWT Settings
```yaml
jwt:
  secret: "your-secret-key"
  expiration: 24  # hours
```

## 📈 Performance & Scalability

- **Connection Pooling**: Configured max connections and idle connections
- **Rate Limiting**: 100 requests per minute per IP
- **Request Size Limit**: 10MB maximum
- **Database Indexes**: Optimized for common queries
- **Caching**: Ready for Redis integration
- **Horizontal Scaling**: Stateless API design

## 🛡️ Security Features

- JWT-based authentication
- Password hashing with bcrypt
- Role-based access control (RBAC)
- Request size validation
- Rate limiting
- CORS configuration
- SQL injection protection (via GORM)
- Environment variable for sensitive data

## 🚀 Deployment

### Deploy to Railway
```bash
# Railway will auto-detect Dockerfile
railway up
```

### Deploy to Azure
1. Push to Azure Container Registry
2. Deploy to Azure App Service or AKS

### Deploy to AWS
1. Push to ECR
2. Deploy to ECS or EKS

## 📝 Database Schema

### Main Tables
- **users**: User accounts and authentication
- **vehicles**: Customer vehicles
- **waiting_lists**: Queue management
- **maintenance_items**: Maintenance tasks and approvals
- **products**: Parts and service inventory
- **parts**: Part details
- **invoices**: Billing information
- **settings**: Application settings

## 🔄 Architecture

This service follows **Clean Architecture** principles:

1. **Domain Layer**: Business entities and rules
2. **Use Case Layer**: Application business logic
3. **Interface Layer**: HTTP handlers and external adapters
4. **Infrastructure Layer**: Database, messaging, logging

### Design Patterns Used
- Repository Pattern
- Dependency Injection
- Factory Pattern
- Strategy Pattern

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👥 Authors

- **kuahbanyak** - *Initial work*

## 🙏 Acknowledgments

- Go community for excellent libraries
- Clean Architecture by Robert C. Martin
- Domain-Driven Design principles

## 📞 Support

For support and questions:
- Create an issue on GitHub
- Email: support@example.com

## 🗺️ Roadmap

- [ ] Add Redis caching
- [ ] Implement WebSocket for real-time updates
- [ ] Add payment gateway integration (Stripe)
- [ ] Implement email notifications
- [ ] Add SMS notifications
- [ ] Mobile app integration
- [ ] Advanced analytics dashboard
- [ ] Multi-language support
- [ ] Export reports (PDF/Excel)
- [ ] Integration with third-party services

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:8080/health
```

### RabbitMQ Management
Access RabbitMQ management UI:
```
http://localhost:15672
Username: admin
Password: rabbitmq_secure_password_123
```

## 🔍 Troubleshooting

### Database Connection Issues
- Verify SQL Server is running
- Check connection string in `.env`
- Ensure database exists
- Verify firewall rules

### RabbitMQ Connection Issues
- Verify RabbitMQ is running
- Check credentials in `.env`
- Ensure port 5672 is accessible

### Docker Issues
```bash
# Clean up and rebuild
docker-compose down -v
docker-compose up --build
```

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [GORM Documentation](https://gorm.io/)
- [Gorilla Mux](https://github.com/gorilla/mux)
- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html)
- [Docker Documentation](https://docs.docker.com/)

---

**Built with ❤️ using Go**

