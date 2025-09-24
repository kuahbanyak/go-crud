-- Create database if it doesn't exist
IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'gocrud')
BEGIN
    CREATE DATABASE gocrud;
END
GO

PRINT 'Database gocrud created successfully!';
