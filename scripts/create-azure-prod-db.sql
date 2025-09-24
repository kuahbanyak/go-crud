-- Run this script on Azure Data Studio or SQL Server Management Studio
-- Connect to: devdbsql.database.windows.net
-- Authentication: SQL Server Authentication
-- Login: devsql
-- Password: Kuahpisah1
-- Database: master

-- Create production database
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

-- Verify database creation
SELECT name FROM sys.databases WHERE name = 'sqlprod';
GO
