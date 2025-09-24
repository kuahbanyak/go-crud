-- Create database if it doesn't exist
IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'gocrud')
BEGIN
    CREATE DATABASE gocrud;
END
GO

USE gocrud;
GO

-- Create Users table
IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='users' AND xtype='U')
CREATE TABLE users (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    username NVARCHAR(50) UNIQUE NOT NULL,
    email NVARCHAR(100) UNIQUE NOT NULL,
    password_hash NVARCHAR(255) NOT NULL,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE()
);
GO

-- Create Products table
IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='products' AND xtype='U')
CREATE TABLE products (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    name NVARCHAR(100) NOT NULL,
    description NVARCHAR(500),
    price DECIMAL(10,2) NOT NULL,
    stock_quantity INT DEFAULT 0,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE()
);
GO

-- Create Vehicles table
IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='vehicles' AND xtype='U')
CREATE TABLE vehicles (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    make NVARCHAR(50) NOT NULL,
    model NVARCHAR(50) NOT NULL,
    year INT NOT NULL,
    license_plate NVARCHAR(20) UNIQUE,
    user_id UNIQUEIDENTIFIER,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
GO

-- Create Bookings table
IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='bookings' AND xtype='U')
CREATE TABLE bookings (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    user_id UNIQUEIDENTIFIER NOT NULL,
    vehicle_id UNIQUEIDENTIFIER,
    service_type NVARCHAR(100) NOT NULL,
    booking_date DATETIME2 NOT NULL,
    status NVARCHAR(20) DEFAULT 'pending',
    total_amount DECIMAL(10,2),
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE(),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id)
);
GO

-- Insert sample data (only if tables are empty)
IF NOT EXISTS (SELECT * FROM users)
BEGIN
    INSERT INTO users (username, email, password_hash) VALUES
    ('admin', 'admin@example.com', '$2a$10$hash_placeholder_for_admin_password'),
    ('testuser', 'test@example.com', '$2a$10$hash_placeholder_for_test_password');
END
GO

IF NOT EXISTS (SELECT * FROM products)
BEGIN
    INSERT INTO products (name, description, price, stock_quantity) VALUES
    ('Engine Oil', 'High quality synthetic engine oil', 29.99, 100),
    ('Brake Pads', 'Premium brake pads for cars', 79.99, 50),
    ('Air Filter', 'Replacement air filter', 19.99, 75);
END
GO

PRINT 'Database initialization completed successfully!';
