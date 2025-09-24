# Azure SQL Database Production Setup Script

# Connect to your Azure SQL Server and create production database
# Run this in Azure Data Studio or SQL Server Management Studio

-- Connect to: devdbsql.database.windows.net
-- User: devsql
-- Password: Kuahpisah1
-- Database: master (to create new database)

-- Create production database
IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'sqlprod')
BEGIN
    CREATE DATABASE sqlprod;
END
GO

PRINT 'Production database sqlprod created successfully on devdbsql.database.windows.net';

-- Switch to the new database
USE sqlprod;
GO

-- Grant permissions to your existing user (if needed)
-- The user 'devsql' should already have access since it's on the same server

PRINT 'Production database setup completed!';
PRINT 'Your Go application will auto-create tables using GORM migration on first startup.';
