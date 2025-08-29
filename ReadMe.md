# Car Service Management Backend

A comprehensive Go-based REST API for managing car service operations, built with Gin framework and SQL Server.

## 🚀 Features

- **User Management**: Registration, authentication, JWT-based authorization with role-based access control
- **Vehicle Management**: Complete CRUD operations for customer vehicles
- **Booking System**: Service appointment scheduling and status management
- **Service History**: Detailed service records with receipt upload functionality
- **Inventory Management**: Parts tracking and stock management
- **Invoice Generation**: Automated billing with PDF generation
- **Notifications**: Email notifications for booking updates
- **Reporting**: Business analytics and reporting endpoints

## 🛠 Tech Stack

- **Backend**: Go 1.24+ with Gin web framework
- **Database**: SQL Server with GORM ORM
- **Authentication**: JWT tokens with role-based access
- **File Storage**: Local file system storage
- **Email**: SMTP integration for notifications
- **Testing**: Comprehensive unit tests with testify

## 📋 Prerequisites

- Go 1.24 or higher
- SQL Server instance
- SMTP server (for email notifications)

## 🔧 Installation & Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd go-crud
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Environment Configuration

Copy `.env.example` to `.env` and configure your environment variables:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
PORT=8080
DB_DSN=sqlserver://username:password@server:1433?database=dbname
JWT_SECRET=your_super_secret_jwt_key_here
STORAGE_PATH=./uploads
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASS=your_app_password
```

### 4. Database Setup

Ensure your SQL Server is running and the database exists. The application will auto-migrate tables on startup.

### 5. Create Upload Directory

```bash
mkdir uploads
```

### 6. Run the Application

```bash
go run ./cmd/server
```

The server will start on `http://localhost:8080`

## 🏗 Project Structure

```
├── cmd/
│   └── server/          # Application entry point
├── config/              # Configuration management
├── internal/            # Private application code
│   ├── auth/           # Authentication handlers
│   ├── booking/        # Booking management
│   ├── db/             # Database utilities
│   ├── inventory/      # Inventory management
│   ├── invoice/        # Invoice generation
│   ├── server/         # Server setup and routing
│   ├── servicehistory/ # Service records
│   ├── user/           # User management
│   └── vehicle/        # Vehicle management
├── pkg/                # Shared/public packages
│   ├── middleware/     # HTTP middlewares
│   ├── notification/   # Email notifications
│   └── storage/        # File storage
└── test/               # Unit tests
```

## 🔐 API Authentication

The API uses JWT tokens for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

### User Roles

- **admin**: Full system access
- **mechanic**: Service operations access
- **customer**: Limited access to own data

## 📚 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/register` | Register new user | No |
| POST | `/auth/login` | User login | No |

**Register Request:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe",
  "phone": "1234567890",
  "role": "customer",
  "address": "123 Main St"
}
```

**Login Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Users

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/users` | List all users | Admin |
| GET | `/users/:id` | Get user by ID | Admin/Owner |
| PUT | `/users/:id` | Update user | Admin/Owner |

### Vehicles

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/vehicles` | Add new vehicle | Customer/Admin |
| GET | `/vehicles` | List user's vehicles | Yes |
| GET | `/vehicles/:id` | Get vehicle details | Owner/Admin |
| PUT | `/vehicles/:id` | Update vehicle | Owner/Admin |
| DELETE | `/vehicles/:id` | Delete vehicle | Owner/Admin |

**Vehicle Request:**
```json
{
  "brand": "Toyota",
  "model": "Camry",
  "year": 2020,
  "license_plate": "ABC-123",
  "vin": "1HGCM82633A123456",
  "mileage": 50000
}
```

### Bookings

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/bookings` | Create booking | Customer/Admin |
| GET | `/bookings` | List user's bookings | Yes |
| GET | `/bookings/:id` | Get booking details | Owner/Admin |
| PUT | `/bookings/:id/status` | Update booking status | Mechanic/Admin |

**Booking Request:**
```json
{
  "vehicle_id": 1,
  "scheduled_at": "2024-12-25T10:00:00Z",
  "duration_min": 120,
  "notes": "Regular maintenance"
}
```

**Status Update:**
```json
{
  "status": "in_progress"
}
```

**Booking Statuses:**
- `scheduled` - Appointment scheduled
- `in_progress` - Service in progress
- `completed` - Service completed
- `canceled` - Booking canceled

### Inventory

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/inventory/parts` | Add new part | Admin/Mechanic |
| GET | `/inventory/parts` | List all parts | Admin/Mechanic |
| PUT | `/inventory/parts/:id` | Update part | Admin/Mechanic |

**Part Request:**
```json
{
  "sku": "BRK-001",
  "name": "Brake Pad",
  "qty": 50,
  "price": 25000
}
```

### Service History

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/service-history` | Add service record | Mechanic/Admin |
| GET | `/vehicles/:id/history` | Get vehicle service history | Owner/Admin |

**Service Record Request:**
```json
{
  "booking_id": 1,
  "vehicle_id": 1,
  "description": "Oil change and brake inspection",
  "cost": 150000,
  "receipt_url": "https://example.com/receipt.pdf"
}
```

### Invoices

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/invoices` | Generate invoice | Admin/Mechanic |
| GET | `/invoices/summary` | Get invoice summary | Admin |

**Invoice Request:**
```json
{
  "booking_id": 1,
  "amount": 200000,
  "status": "pending"
}
```

## 🧪 Testing

The project includes comprehensive unit tests for all services.

### Run All Tests

```bash
go test ./test/... -v
```

### Run Individual Test Files

```bash
# Run user tests
go test -v ./test/user_test.go ./test/main_test.go

# Run vehicle tests
go test -v ./test/vehicle_test.go ./test/main_test.go

# Run booking tests
go test -v ./test/booking_test.go ./test/main_test.go
```

### Custom Test Runner

For environments without CGO support:

```bash
go run ./test/runner.go
```

### Test Coverage

- Repository pattern testing with mocks
- Model validation testing
- Error handling scenarios
- Integration workflow testing

## 📊 Database Schema

### Users Table
- `id` - Primary key
- `email` - Unique user email
- `password` - Hashed password
- `name` - Full name
- `phone` - Contact number
- `role` - User role (admin/mechanic/customer)
- `address` - Physical address

### Vehicles Table
- `id` - Primary key
- `owner_id` - Foreign key to users
- `brand` - Vehicle brand
- `model` - Vehicle model
- `year` - Manufacturing year
- `license_plate` - License plate number
- `vin` - Vehicle identification number
- `mileage` - Current mileage

### Bookings Table
- `id` - Primary key
- `vehicle_id` - Foreign key to vehicles
- `customer_id` - Foreign key to users
- `mechanic_id` - Foreign key to users (optional)
- `scheduled_at` - Appointment date/time
- `duration_min` - Service duration in minutes
- `status` - Booking status
- `notes` - Additional notes

### Parts Table (Inventory)
- `id` - Primary key
- `sku` - Stock keeping unit
- `name` - Part name
- `qty` - Quantity in stock
- `price` - Part price

### Service Records Table
- `id` - Primary key
- `booking_id` - Foreign key to bookings
- `vehicle_id` - Foreign key to vehicles
- `description` - Service description
- `cost` - Service cost
- `receipt_url` - Receipt file URL

### Invoices Table
- `id` - Primary key
- `booking_id` - Foreign key to bookings
- `amount` - Invoice amount
- `status` - Payment status
- `pdf_url` - Generated PDF URL

## 🚀 Deployment

### Building for Production

```bash
go build -o car-service-api ./cmd/server
```

### Docker Deployment

Create a `Dockerfile`:

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./main"]
```

Build and run:

```bash
docker build -t car-service-api .
docker run -p 8080:8080 car-service-api
```

## 📝 Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_DSN` | Database connection string | `sqlserver://user:pass@host:1433?database=db` |
| `JWT_SECRET` | JWT signing secret | `your_secret_key` |
| `STORAGE_PATH` | File upload directory | `./uploads` |
| `SMTP_HOST` | SMTP server host | `smtp.gmail.com` |
| `SMTP_PORT` | SMTP server port | `587` |
| `SMTP_USER` | SMTP username | `your_email@gmail.com` |
| `SMTP_PASS` | SMTP password | `your_app_password` |

## 🔧 Configuration

The application uses GORM for database operations with auto-migration enabled. On first run, it will create all necessary tables.

### CORS Configuration

The API includes CORS middleware configured to allow requests from any origin in development. Update for production use.

### File Upload

Files are stored locally in the `STORAGE_PATH` directory. Implement cloud storage (AWS S3, etc.) for production.

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Verify SQL Server is running
   - Check connection string format
   - Ensure database exists

2. **JWT Token Invalid**
   - Check JWT_SECRET configuration
   - Verify token format in Authorization header

3. **File Upload Errors**
   - Ensure STORAGE_PATH directory exists
   - Check directory permissions

4. **Email Notifications Not Working**
   - Verify SMTP configuration
   - Check firewall settings
   - Use app passwords for Gmail

### Logging

The application uses Gin's default logging. Enable debug mode for detailed logs:

```bash
GIN_MODE=debug go run ./cmd/server
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the troubleshooting section
- Review the API documentation

## 🔄 Changelog

### v1.0.0
- Initial release
- Complete CRUD operations for all entities
- JWT authentication with role-based access
- File upload functionality
- Email notifications
- Comprehensive test suite

---

**Note**: This is a development scaffold. Enhance security measures, implement proper error handling, and add monitoring for production deployment.
