# 🚗 Vehicle Service Management System

A comprehensive backend API for managing vehicle service operations, built with Go, Gin, and GORM. This system provides complete functionality for auto repair shops, service centers, and vehicle maintenance businesses.

## 🌟 Features

### Core Functionality
- **👥 User Management** - Registration, authentication, JWT-based authorization
- **🚙 Vehicle Management** - Complete CRUD operations for customer vehicles
- **📅 Booking System** - Service appointment scheduling with status tracking
- **📋 Service History** - Detailed service records and maintenance history
- **📦 Inventory Management** - Parts tracking, stock management
- **💰 Invoice Generation** - Automated billing and PDF invoice generation
- **📧 Notifications** - Email notifications and real-time alerts

### Advanced Features
- **💬 Real-time Communication** - WebSocket-based messaging between customers and mechanics
- **📆 Advanced Scheduling** - Mechanic availability, service types, maintenance reminders
- **🎯 Service Packages** - Predefined service bundles and categories
- **📊 Customer Dashboard** - Vehicle health monitoring, service recommendations
- **⏰ Maintenance Reminders** - Automated reminder system based on time/mileage
- **📋 Waitlist Management** - Customer waitlist for busy periods

## 🏗️ Architecture

```
go-crud/
├── cmd/server/          # Application entry point
├── config/              # Configuration management
├── internal/            # Private application code
│   ├── auth/           # Authentication & authorization
│   ├── booking/        # Booking management
│   ├── dashboard/      # Customer dashboard
│   ├── inventory/      # Parts & inventory
│   ├── invoice/        # Billing & invoicing
│   ├── message/        # Real-time messaging
│   ├── notification/   # WebSocket notifications
│   ├── scheduling/     # Advanced scheduling
│   ├── servicehistory/ # Service records
│   ├── servicepackage/ # Service packages
│   ├── user/           # User management
│   ├── vehicle/        # Vehicle management
│   └── server/         # HTTP server & routing
└── pkg/                # Public packages
    ├── middleware/     # HTTP middleware
    ├── notification/   # Email notifications
    └── storage/        # File storage
```

## 🚀 Quick Start

### Prerequisites
- Go 1.19+
- PostgreSQL or SQL Server
- SMTP server (for email notifications)

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd go-crud
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Set up environment variables**
Create a `.env` file in the root directory:
```env
PORT=8080
DB_DSN=sqlserver://username:password@server:1433?database=dbname
JWT_SECRET=your-super-secret-jwt-key
STORAGE_PATH=./uploads
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
```

4. **Run the application**
```bash
go run ./cmd/server
```

The server will start on `http://localhost:8080`

## 📡 API Documentation

### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - User login

### Core Operations
- `GET|POST|PUT|DELETE /api/v1/vehicles` - Vehicle management
- `GET|POST|PUT /api/v1/bookings` - Booking management  
- `GET|POST /api/v1/service-records` - Service history
- `GET|POST|PUT /api/v1/parts` - Inventory management
- `POST /api/v1/invoices/generate` - Generate invoices

### Advanced Features
- `GET /api/v1/ws` - WebSocket connection for real-time updates
- `POST /api/v1/messages` - Send messages
- `GET|POST /api/v1/scheduling/*` - Scheduling operations
- `GET|POST /api/v1/packages/*` - Service packages
- `GET /api/v1/dashboard/*` - Dashboard analytics

### 📥 Postman Collection
Import the included `go-crud-api.postman_collection.json` file into Postman for complete API testing with:
- Pre-configured environment variables
- Automatic JWT token handling
- Sample requests for all endpoints
- Organized folder structure

## 🛠️ Development

### Database Migration
The application automatically handles database migrations on startup using GORM's AutoMigrate feature.

### Building
```bash
# Build for current platform
go build -o server ./cmd/server

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o server.exe ./cmd/server

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o server ./cmd/server
```

### Testing
```bash
go test ./...
```

## 🌐 WebSocket Integration

The system supports real-time communication through WebSocket connections:

```javascript
// Client-side WebSocket connection
const ws = new WebSocket('ws://localhost:8080/api/v1/ws');
ws.onmessage = function(event) {
    const notification = JSON.parse(event.data);
    // Handle real-time notifications
};
```

## 📊 Key Models

- **User** - Customer and mechanic profiles
- **Vehicle** - Vehicle information and ownership
- **Booking** - Service appointments and scheduling
- **ServiceHistory** - Completed service records
- **Inventory** - Parts and stock management
- **Invoice** - Billing and payment tracking
- **Message** - Real-time communication
- **ServicePackage** - Predefined service bundles

## 🔒 Security Features

- JWT-based authentication
- CORS middleware
- Request rate limiting
- SQL injection protection via GORM
- Input validation and sanitization

## 🚀 Deployment

### Docker (Recommended)
```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### Environment Configuration
- **Development**: Use `.env` file
- **Production**: Set environment variables directly
- **Docker**: Use docker-compose with env files

## 📈 Monitoring & Health Checks

The API includes basic health monitoring endpoints and structured logging for production deployments.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📝 License

This project is licensed under the MIT License.

## 📞 Support

For support and questions, please create an issue in the repository.

---

**Note**: This is a production-ready scaffold. Ensure you update security configurations, environment variables, and add proper error handling before deploying to production.
