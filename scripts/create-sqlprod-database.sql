-- Connect to Azure SQL Database using Azure Data Studio or SQL Server Management Studio
-- Server: devdbsql.database.windows.net
-- Authentication: SQL Server Authentication
-- Login: devsql
-- Password: Kuahpisah1
-- Database: master (connect to master first to create new database)

-- Create the production database
IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'sqlprod')
BEGIN
    CREATE DATABASE sqlprod;
    PRINT 'Production database sqlprod created successfully!';
END
ELSE
BEGIN
    PRINT 'Production database sqlprod already exists.';
END
GO

-- Switch to the new production database context
USE sqlprod;
GO

-- Verify we're connected to the right database
SELECT DB_NAME() AS CurrentDatabase;
GO

PRINT 'Ready for production deployment!';
GO
